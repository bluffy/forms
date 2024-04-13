package router

import (
	"io/fs"
	"net/http"
	"strings"
	"sync/atomic"

	"goapp/app"
	"goapp/app/middleware"

	chi_middleware "github.com/go-chi/chi/v5/middleware"

	"gitea.com/go-chi/session"
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
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
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
	r.Use(chi_middleware.Heartbeat("/ping"))
	r.Use(chi_middleware.RealIP)

	//r.Use(session.Sessioner())
	r.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "tmp/sessions",
		CookieName:     "session",
		IDLength:       64,
	}))
	r.Use(middleware.Logger("", nil))
	r.Use(cors.Handler(cors.Options{
		//AllowedOrigins:   []string{"*"},
		//AllowedOrigins:   []string{"http://localhost*", "http://127.0.0.1*", "http://127.0.0.1*", "http://128.140.68.242"},
		AllowedOrigins:   a.Conf().Server.Cors.AllowedOrigins,
		AllowCredentials: a.Conf().Server.Cors.AllowCredentials,
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: a.Conf().Server.Cors.AllowedMethods,
		AllowedHeaders: a.Conf().Server.Cors.AllowedHeaders,
		ExposedHeaders: a.Conf().Server.Cors.ExposedHeaders,

		MaxAge: a.Conf().Server.Cors.MaxAge, // Maximum value not ignored by any of major browsers
	}))
	//r.Use(middleware.SetSession(a))
	r.Use(middleware.SetLocale(a))

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
	//r.Get("/home", a.PageHome)

	fileServer(r, "/pdf", http.FS(html2pdfFS))

	/*
		r.Get("/test", func(sess session.Store) string {
			sess.Set("session", "session middleware")
			return sess.Get("session").(string)
		})
		r.Get("/test1", func(w http.ResponseWriter, r *http.Request, f *session.Flash) string {
			f.Success("yes!!!")
			f.Error("opps...")
			f.Info("aha?!")
			f.Warning("Just be careful.")
			w.HTML(200, "signup")
		})
	*/

	r.Route("/p", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			//		r.Get("/register/{link}", a.HandlerRegisterLinkGet)
		})
	})

	r.Route("/bl-api", func(r chi.Router) {
		if a.Conf().Dev || a.Conf().ShowApiDoku {
			r.Mount("/doc", httpSwagger.WrapHandler)
		}
		r.Route("/v1", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.ContentTypeJson)
				r.Use(middleware.CheckUserLogin(a))
				r.Get("/", a.PageIndex)

			})
			r.Route("/user", func(r chi.Router) {

				r.Group(func(r chi.Router) {
					r.Use(middleware.ContentTypeJson)
					r.Post("/login", a.HandleAuthLoginForm)
					r.Post("/register/link", a.HanderCreateUserFromMailLink)
					r.Post("/register", a.HandleAuthRegisterFrom)
					r.Post("/forgot_password/link", a.HandlerCreateNewPasswordFromMailLink)
					r.Post("/forgot_password", a.HandlerGenerateMailWithPasswordLink)
				})
				r.Group(func(r chi.Router) {
					r.Use(middleware.ContentTypeJson)
					r.Use(middleware.CheckUserLogin(a))
					r.Get("/", a.HanldeCheckUser)
					r.Get("/logout", a.HandleAuthLogout)

				})

			})

		})

		//r.Get("/oidc/{name}", a.HandlerOpenIDConnect)

		/*
			r.Group(func(r chi.Router) {
				r.Use(middleware.ContentTypeJson)
				r.Use(middleware.SessionCheck(a))
				r.Get("/", a.PageIndex)
			})
		*/
		r.Group(func(r chi.Router) {
			/*
				r.Use(middleware.ContentTypeJson)

				r.Post("/login", a.HandlerLogin)
				r.Post("/register", a.HandlerRgister)
			*/
			//r.Post("/oidc/callback/{name}", a.HandlerOpenIDCallback)
			/*
						r.Get("/login", a.HandlerLoginData)

				r.Get("/forgot_password", a.HandlePasswordRequest)
				r.Post("/login/refresh", a.RefreshLoginToken)

				r.Post("/user", a.HandlerLogin)
				r.Post("/new_password", a.HandlerNewPassword)
			*/
			/*
				r.Group(func(r chi.Router) {
					r.Use(middleware.JWTAuth(a))
					//	r.Get("/user", a.HandlerSessionCheck)
					//	r.Get("/user", a.HandlerSessionCheck)
					//		r.Delete("/user", a.HandlerSessionRemove)
				})*/
			//r.HandleFunc("/{user:user\\/?}", userLogin).Methods("POST")
			//r.HandleFunc("/{user:user\\/?}", userLogin).Methods("POST")
			//	r.Handle("/{user:user\\/?}", auth.AuthMiddleware(userGet)).Methods("GET")
			//		r.Handle("/{user:user\\/?}", auth.AuthMiddleware(userDelete)).Methods("DELETE")
			//r.Post("/auth/refresh", a.RefreshLoginToken)
			/*
				r.Group(func(r chi.Router) {
					//r.Use(middleware.JWTAuth(a))
					//r.Use(middleware.SessionCheck(a))
					r.Get("/intern", a.HandlerIntern)
					//	r.Get("/user", a.HandlerSessionCheck)
					//	r.Get("/user", a.HandlerSessionCheck)
					//		r.Delete("/user", a.HandlerSessionRemove)
				})*/
		})

		//sicherheit muss noch eingebuat werden
		//r.Get("/intern", a.HanldeIntern)
		/*
			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(a))
				r.Use(middleware.ContentTypeJson)
				//, auth.AuthMiddleware(FileGetHtmlPdfAuth)).Methods("GET")

			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.JWTAuth(a))
				r.Use(middleware.ContentTypeJson)
			})
		*/
	})
	/*
		r.Get("/", a.PageHome)
	*/

	if a.Conf().UseEmbedClient {
		r.Get("/*", a.HandleClient)
	} else {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})
	}

	return r

}
