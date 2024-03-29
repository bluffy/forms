package router

import (
	"blubooks/server/app"
	"blubooks/server/app/middleware"

	//"blubooks/server/app/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func New(a *app.App) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger("", nil))
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		//AllowedOrigins: []string{"https://*", "http://*"},
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Get("/healthz", app.HandleHealth)
	//r.Use(middleware.Logger("", nil))
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.ContentTypeJson)
			r.Post("/auth/login", a.Login)
			r.Post("/auth/refresh", a.RefreshLoginToken)
		})

		r.Route("/page", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(a))
				r.Use(middleware.ContentTypeJson)
				//r.Get("/client/{id}/collections", a.PageGetCollectionsFromClient)
				r.Get("/client/{id}", a.PageReadClient)
				r.Get("/collection/{id}", a.PageGetCollection)
				r.Get("/section/{id}", a.PageReadSection)
				r.Get("/book/{id}", a.PageReadBook)
				//r.Get("/clients/{id}/collections", a.PageGetCollections)
				//r.Get("/clients/{id}/collections", a.GetCollections)

				r.Get("/clients", a.PageListClients)
			})

		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(a))
			r.Use(middleware.ContentTypeJson)
			r.Get("/client/{id}", a.ReadClient)
			r.Get("/section/{id}", a.ReadSection)
			r.Put("/section/{id}", a.UpdateSection)
			r.Post("/section/{id}", a.CreateSection)
			//r.Get("/clients/{id}/collections", a.GetCollections)
			//r.Get("/clients/{id}/collections", a.GetCollections)

			r.Get("/client", a.ListClients)
		})
	})
	r.Get("/*", a.HandleIndex)
	return r
}
