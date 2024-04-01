package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/segmentio/ksuid"

	dbConn "goapp/adapter/gorm"
	"goapp/app"
	"goapp/config"
	"goapp/lang"
	"goapp/router"
	"goapp/util/logger/goose_logger"
	"goapp/util/tools"

	goose_v3 "github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"

	_ "goapp/docs" // swagger
	_ "goapp/migrations"

	"gorm.io/gorm"

	vr "goapp/util/validator"
)

type ArgOptions struct {
	Config  string `short:"c" long:"config" description:"config.yaml file"`
	Migrate string `short:"m" long:"migrate" description:"DB mirgrate tool" choice:"up" choice:"down" choice:"status" choice:"version" choice:"reset" choice:"up-by-one" choice:"up-to" choice:"down-to"`
	PWHash  string `short:"p" long:"password" description:"Password Hash"`
	UID     bool   `short:"u" long:"uid" description:"UID"`
}

/*
type ArgOptions struct {
	Env string `short:"e" long:"env" description:"Environment File (default: .env)"`
	//Migrate string `short:"m" long:"migrate" description:"DB mirgrate tool"  choice:"up" choice:"down" choice:"status" choice:"version" choice:"reset" choice:"up-by-one" choice:"up-to" choice:"down-to"`
	//PWHash  string `short:"p" long:"password" description:"Password Hash"`
}
*/

//go:embed migrations
var embedMigrations embed.FS

//go:embed public/*
var publicFS embed.FS

//go:embed data/*
var dataFS embed.FS

func main() {
	var opts ArgOptions
	var err error
	var args []string
	var configFile = "config.yaml"

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	args, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}

	if opts.Config != "" {
		configFile = opts.Config
	}

	_, err = config.LoadConfig(configFile)
	if err != nil {
		log.WithField("error", err).Fatal("Error in reading Config File")
		return
	}

	if config.Conf.Debug {
		log.Info("DEBUG Mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.Info("PRODUCTION Mode")
		log.SetLevel(log.InfoLevel)
	}

	if opts.PWHash != "" {
		hashed, err := tools.HashPassword(opts.PWHash)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Password: %s Hashed: %s", opts.PWHash, hashed)

		return
	}
	if opts.UID {
		uid := ksuid.New().String()
		log.Printf("UID: %s", uid)
		return
	}

	Server(opts, args)
}

// @title  app server
// @version 1.0
// @description app server

// @contact.name API Support
// @contact.email github@bluffy.de

// @schemes  http https
// @BasePath /

// @securityDefinitions.apikey BEARER
// @in header
// @name Authorization
// @description Type "Token" followed by a space and JWT token.
// Server function creates start Listenen Server
func Server(opts ArgOptions, args []string) {

	var db *gorm.DB
	var err error

	log.Info("Set Language: " + config.Conf.Language)

	appLang := lang.AppLang(dataFS)

	// check DB Connection on Start 100 times
	log.Info("Connect Database: " + config.Conf.Database.Type)

	for i := 1; i <= 100; i++ {
		db, err = dbConn.New()
		if err != nil {
			log.Error(err)
		} else {
			break
		}
		time.Sleep(4 * time.Second)
	}
	if err != nil {
		log.Error("Start Failed: Stopped Trying DB Connection")
		return
	}
	log.Info("Start Migrate: " + config.Conf.Database.Type)
	if migrate(db, opts.Migrate, args, config.Conf.Database.Type) {
		log.Info("Program exited")
		return
	}
	log.Info("Server Start")
	/*
		err = db.Ping()
		if err != nil {
			log.Error("Start Failed: Databse not ready")
			return
		}*/

	validator := vr.New()

	addressApp := fmt.Sprintf(":%d", config.Conf.Server.Port)
	addressApi := fmt.Sprintf(":%d", config.Conf.Server.PortIntern)

	application := app.New(validator, appLang, db)

	appRouter := router.NewApp(application, publicFS)
	internRouter := router.NewIntern(application, publicFS)

	srv := &http.Server{
		Addr:         addressApp,
		Handler:      appRouter,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	srvIntern := &http.Server{
		Addr:         addressApi,
		Handler:      internRouter,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	log.Infof("Starting APP Server %v", addressApp)
	go srv.ListenAndServe()
	log.Infof("Starting INTERN Server %v", addressApi)
	go srvIntern.ListenAndServe()

	/*
		if config.Conf.Debug == true {
			routerLiveReload := http.NewServeMux()
			logger := logger.NewLogger(os.Stdout, logger.LogLevelInfo, true)
			// initialize http.Server
			srvLiveReload := &http.Server{
				Addr:    ":35729",
				Handler: routerLiveReload,
			}

			// start the application with live reload
			gloader.NewWatcher(srvLiveReload, time.Second*2, logger).Start()
		}
	*/

	// create a channel to subscribe ctrl+c/SIGINT event
	sigInterruptChannel := make(chan os.Signal, 1)

	signal.Notify(sigInterruptChannel, os.Interrupt)
	// block execution from continuing further until SIGINT comes
	<-sigInterruptChannel

	// create a context which will expire after 4 seconds of grace period
	timeout := time.Second * 4
	if config.Conf.Dev {
		timeout = time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Info("Shutting down APP Server")
	go srv.Shutdown(ctx)
	log.Info("Shutting down API Server")

	go srvIntern.Shutdown(ctx)

	sqlDB, err := db.DB()
	if err == nil {
		log.Info("Shutting down DB Server")
		if err = sqlDB.Close(); err != nil {
			log.WithField("error", err).Warn("Db connection closing failure")
		}
	}

	// wait until ctx ends (which will happen after 4 seconds)
	<-ctx.Done()

}

func migrate(db *gorm.DB, migrateCMD string, migrateArgs []string, dialect string) (doExit bool) {
	appDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return
	}

	goose_v3.SetBaseFS(embedMigrations)

	goose_v3.SetLogger(goose_logger.New())

	if err := goose_v3.SetDialect(dialect); err != nil {
		log.Fatal(err)
	}

	if migrateCMD == "up" {
		if err := goose_v3.Up(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "down" {
		if err := goose_v3.Down(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "status" {
		if err := goose_v3.Status(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "version" {
		if err := goose_v3.Version(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "reset" {
		if err := goose_v3.Reset(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "up-by-one" {
		if err := goose_v3.Reset(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "up-to" || migrateCMD == "down-to" {
		if len(migrateArgs) < 2 {
			log.Fatal("Version missing")
			return true
		}
		version, err := strconv.ParseInt(migrateArgs[1], 10, 64)
		if err != nil {
			log.WithError(err).Fatal("incorrect verison format")
		}

		if migrateCMD == "up-to" {
			if err := goose_v3.UpTo(appDb, "migrations", version); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := goose_v3.DownTo(appDb, "migrations", version); err != nil {
				log.Fatal(err)
			}
		}

	}
	if migrateCMD != "" {
		return true
	}

	if err := goose_v3.Up(appDb, "migrations"); err != nil {
		log.Fatal(err)
	}

	if err := goose_v3.Status(appDb, "migrations"); err != nil {
		log.Fatal(err)
	}

	return false
}
