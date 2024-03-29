package main

import (
	"github.com/bluffy/forms/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	Server()
}

// @title qProCheck API Server
// @version 1.0
// @description QPC

// @contact.name API Support
// @contact.email mritter@ci-database.de

// @schemes  http https
// @BasePath /

// @securityDefinitions.apikey Token
// @in header
// @name Authorization
// @description Type "Token" followed by a space and JWT token.
// Server function creates start Listenen Server
func Server() {

	log.Info("Server Start")

	appConf, err := config.AppConfig("config.yaml")
	if err != nil {
		log.WithField("error", err).Fatal("Error in reading Config File")
	}

	if appConf.Debug {
		log.Info("DEBUG Mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.Info("PRODUCTION Mode")
	}

	/*

		var db *sql.DB
		// check DB Connection on Start 100 times
		for i := 1; i <= 10000; i++ {
			db, err = dbConn.Open(appConf)
			if err != nil {
				log.WithField("error", err).Warn("Could not Connect to database")
			} else {
				break
			}
			time.Sleep(4 * time.Second)

		}
		if err != nil {
			log.Error("Start Failed: Stopped Trying DB Connection")
			return
		}

		err = db.Ping()
		if err != nil {
			log.Error("Start Failed: Databse not ready")
			return
		}

		addressApp := fmt.Sprintf(":%d", appConf.Server.PortApp)
		addressApi := fmt.Sprintf(":%d", appConf.Server.PortApi)

		application := app.New(appConf, appLang, db)
		appRouter := router.NewApp(application, publicFS, dataFS)
		apiRouter := router.NewApi(application, publicFS, dataFS)
		//	appRouter.Mount("/swagger", httpSwagger.WrapHandler)
		//	apiRouter.Mount("/swagger", httpSwagger.WrapHandler)

		srv := &http.Server{
			Addr:         addressApp,
			Handler:      appRouter,
			ReadTimeout:  appConf.Server.TimeoutRead,
			WriteTimeout: appConf.Server.TimeoutWrite,
			IdleTimeout:  appConf.Server.TimeoutIdle,
		}

		srvApi := &http.Server{
			Addr:         addressApi,
			Handler:      apiRouter,
			ReadTimeout:  appConf.Server.TimeoutRead,
			WriteTimeout: appConf.Server.TimeoutWrite,
			IdleTimeout:  appConf.Server.TimeoutIdle,
		}

		log.Infof("Starting APP Server %v", addressApp)
		go srv.ListenAndServe()
		log.Infof("Starting API Server %v", addressApi)
		go srvApi.ListenAndServe()

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

		go srvApi.Shutdown(ctx)

		if err = dbConn.Close(); err != nil {
			log.WithField("error", err).Warn("Db connection closing failure")
		}

		// wait until ctx ends (which will happen after 4 seconds)
		<-ctx.Done()
	*/

}
