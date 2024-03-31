package router

import (
	"io/fs"
	"sync/atomic"
	"time"

	"github.com/bluffy/forms/server/app"
	"github.com/bluffy/forms/server/app/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title bluffy-forms app intern server
// @version 1.0
// @description

// @contact.name API Support
// @contact.email mario@bluffy.de

// @schemes  http https
// @BasePath /

// @securityDefinitions.apikey Token
// @in header
// @name Authorization
// @description Type "Token" followed by a space and JWT token.
func NewIntern(a *app.App, publicFS fs.FS) *chi.Mux {
	isReady := &atomic.Value{}
	isReady.Store(false)

	r := chi.NewRouter()

	r.Use(middleware.Logger("", nil))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.HandleFunc("/readyz", readyz(isReady))
	r.HandleFunc("/healthz", a.HanlderHealth)

	go func() {
		//time.Sleep(10 * time.Second)
		for {
			/*
				errPing := oracle.Ping()
				if errPing != nil {

					isReady.Store(true)
					logrus.Printf("Readyz true")
					break
				}
			*/

			isReady.Store(true)
			logrus.Printf("Readyz true")

			time.Sleep(2 * time.Second)
			break

		}

	}()

	r.Mount("/swagger", httpSwagger.Handler(httpSwagger.InstanceName("intern")))

	return r

}
