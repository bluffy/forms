package router

import (
	"io/fs"
	"net/http"
	"strings"
	"sync/atomic"

	"goapp/app"
	"goapp/app/middleware"
	"goapp/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger"
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

/*

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// まずはリクエストされた通りにファイルを探索
	err := tryRead(assets, "frontend/dist", r.URL.Path, w)
	if err == nil {
		return
	}
	// 見つからなければindex.htmlを返す
	err = tryRead(assets, "frontend/dist", "index.html", w)
	if err != nil {
		panic(err)
	}
}
*/

func readyz(isReady *atomic.Value) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if isReady == nil || !isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func NewApp(a *app.App, publicFS fs.FS) *chi.Mux {

	isReady := &atomic.Value{}
	isReady.Store(false)

	/*
		go func() {
			logrus.Printf("Readyz probe is negative by default...")
			time.Sleep(10 * time.Second)
			isReady.Store(true)
			logrus.Printf("Readyz probe is positive.")
		}()
	*/

	r := chi.NewRouter()
	r.Use(middleware.Logger("", nil))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.HandleFunc("/healthz", a.HanlderHealth)

	/*
		r.HandleFunc("/readyz", readyz(isReady))
		r.HandleFunc("/healthz", a.HanlderHealth)

	*/

	/*	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})
	*/

	//workDir, _ := os.Getwd()
	//assetsDir := http.Dir("./public/assets/")
	//pdfDir := http.Dir(config.Conf.File.FileDataDir)

	//fileServer(r, "/assets", assetsDir)
	//pdfDir := http.Dir("./public/assets/")
	//fileServer(r, "/html2pdf", pdfDir)

	//publicDirHandler := http.StripPrefix("public", http.FileServer(http.FS(publicFS)))
	//r.Handle("/public/assets/*", publicDirHandler)

	staticFS, _ := fs.Sub(publicFS, "public/static")
	html2pdfFS, _ := fs.Sub(publicFS, "data/html2pdf")

	fileServer(r, "/static", http.FS(staticFS))
	fileServer(r, "/pdf", http.FS(html2pdfFS))
	//fileServerEmbed(r, "/public", publicFS)
	r.Get("/home", a.PageHome)

	r.Route("/api/v1", func(r chi.Router) {

		//r.Get("/oidc/{name}", a.HandlerOpenIDConnect)

		r.Group(func(r chi.Router) {
			r.Use(middleware.ContentTypeJson)

			r.Get("/test", a.HandlerIndex)

			r.Post("/login", a.HandlerLogin)

			//r.Post("/oidc/callback/{name}", a.HandlerOpenIDCallback)
			/*
						r.Get("/login", a.HandlerLoginData)

				r.Get("/forgot_password", a.HandlePasswordRequest)
				r.Post("/login/refresh", a.RefreshLoginToken)

				r.Post("/user", a.HandlerLogin)
				r.Post("/new_password", a.HandlerNewPassword)
			*/
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(a))
				//	r.Get("/user", a.HandlerSessionCheck)
				//	r.Get("/user", a.HandlerSessionCheck)
				//		r.Delete("/user", a.HandlerSessionRemove)
			})
			//r.HandleFunc("/{user:user\\/?}", userLogin).Methods("POST")
			//r.HandleFunc("/{user:user\\/?}", userLogin).Methods("POST")
			//	r.Handle("/{user:user\\/?}", auth.AuthMiddleware(userGet)).Methods("GET")
			//		r.Handle("/{user:user\\/?}", auth.AuthMiddleware(userDelete)).Methods("DELETE")
			//r.Post("/auth/refresh", a.RefreshLoginToken)
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(a))
				r.Get("/intern", a.HandlerIntern)
				//	r.Get("/user", a.HandlerSessionCheck)
				//	r.Get("/user", a.HandlerSessionCheck)
				//		r.Delete("/user", a.HandlerSessionRemove)
			})
		})

		//sicherheit muss noch eingebuat werden
		//r.Get("/intern", a.HanldeIntern)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(a))
			r.Use(middleware.ContentTypeJson)
			//, auth.AuthMiddleware(FileGetHtmlPdfAuth)).Methods("GET")

		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(a))
			r.Use(middleware.ContentTypeJson)
		})
	})
	/*
		r.Get("/", a.PageHome)
	*/

	if config.Conf.UseEmbedClient {
		r.Get("/*", a.HandleClient)
	} else {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})
	}

	r.Mount("/swagger", httpSwagger.WrapHandler)
	return r

}
