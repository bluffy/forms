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

	dbConn "github.com/bluffy/forms/adapter/gorm"
	"github.com/bluffy/forms/app"
	"github.com/bluffy/forms/config"
	"github.com/bluffy/forms/lang"
	"github.com/bluffy/forms/router"
	"github.com/bluffy/forms/util/logger/goose_logger"
	"github.com/bluffy/forms/util/tools"
	goose_v3 "github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"

	_ "github.com/bluffy/forms/docs" // swagger
	"gorm.io/gorm"

	vr "github.com/bluffy/forms/util/validator"
)

/*
type ArgOptions struct {
	Env string `short:"e" long:"env" description:"Environment File (default: .env)"`
	//Migrate string `short:"m" long:"migrate" description:"DB mirgrate tool"  choice:"up" choice:"down" choice:"status" choice:"version" choice:"reset" choice:"up-by-one" choice:"up-to" choice:"down-to"`
	//PWHash  string `short:"p" long:"password" description:"Password Hash"`
}
*/

//go:embed migrations
var embedMigrations embed.FS

type ArgOptions struct {
	Migrate string `short:"m" long:"migrate" description:"DB mirgrate tool" choice:"up" choice:"down" choice:"status" choice:"version" choice:"reset" choice:"up-by-one" choice:"up-to" choice:"down-to"`
	PWHash  string `short:"p" long:"password" description:"Password Hash"`
	UID     bool   `short:"u" long:"uid" description:"UID"`
}

//go:embed public/*
var publicFS embed.FS

//go:embed data/*
var dataFS embed.FS

func main() {
	var opts ArgOptions
	var err error
	var args []string

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	appConf, err := config.AppConfig("config.yaml")
	if err != nil {
		log.WithField("error", err).Fatal("Error in reading Config File")
		return
	}

	if appConf.Debug {
		log.Info("DEBUG Mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.Info("PRODUCTION Mode")
		log.SetLevel(log.InfoLevel)
	}

	args, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
		return
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

	Server(appConf, opts, args)
}

// @title bluffy-forms app server
// @version 1.0
// @description app server

// @contact.name API Support
// @contact.email mario@bluffy.de

// @schemes  http https
// @BasePath /

// @securityDefinitions.apikey BEARER
// @in header
// @name Authorization
// @description Type "Token" followed by a space and JWT token.
// Server function creates start Listenen Server
func Server(appConf *config.Config, opts ArgOptions, args []string) {

	var db *gorm.DB
	var err error

	log.Info("Set Language: " + appConf.Language)

	appLang := lang.AppLang(appConf.Language, appConf.LogLanguage, dataFS)

	// check DB Connection on Start 100 times
	log.Info("Connect Database: " + appConf.Database.Type)

	for i := 1; i <= 100; i++ {
		db, err = dbConn.New(appConf)
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
	log.Info("Start Migrate: " + appConf.Database.Type)
	if migrate(db, opts.Migrate, args, appConf.Database.Type) {
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

	addressApp := fmt.Sprintf(":%d", appConf.Server.Port)
	addressApi := fmt.Sprintf(":%d", appConf.Server.PortIntern)

	application := app.New(appConf, validator, appLang, db)

	appRouter := router.NewApp(application, publicFS)
	internRouter := router.NewIntern(application, publicFS)

	srv := &http.Server{
		Addr:         addressApp,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	srvIntern := &http.Server{
		Addr:         addressApi,
		Handler:      internRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	log.Infof("Starting APP Server %v", addressApp)
	go srv.ListenAndServe()
	log.Infof("Starting INTERN Server %v", addressApi)
	go srvIntern.ListenAndServe()

	/*
		if appConf.Debug == true {
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
	if appConf.Dev {
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
		if err := goose_v3.Up(appDb, "migrations/"+dialect); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "down" {
		if err := goose_v3.Down(appDb, "migrations/"+dialect); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "status" {
		if err := goose_v3.Status(appDb, "migrations/"+dialect); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "version" {
		if err := goose_v3.Version(appDb, "migrations/"+dialect); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "reset" {
		if err := goose_v3.Reset(appDb, "migrations/"+dialect); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "up-by-one" {
		if err := goose_v3.Reset(appDb, "migrations/"+dialect); err != nil {
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
			if err := goose_v3.UpTo(appDb, "migrations/"+dialect, version); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := goose_v3.DownTo(appDb, "migrations/"+dialect, version); err != nil {
				log.Fatal(err)
			}
		}

	}
	if migrateCMD != "" {
		return true
	}

	if err := goose_v3.Up(appDb, "migrations/"+dialect); err != nil {
		log.Fatal(err)
	}

	if err := goose_v3.Status(appDb, "migrations/"+dialect); err != nil {
		log.Fatal(err)
	}

	return false
}
