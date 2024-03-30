package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bluffy/forms/server/app"
	"github.com/bluffy/forms/server/service"
	log "github.com/sirupsen/logrus"
)

func JWTAuth(a *app.App) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			/*
				if r.URL.Query().Get("sess") != "" {
					user := models.User{}
					//user.Token = r.URL.Query().Get("sess")
					a.SetUser(user)
					next.ServeHTTP(w, r)
					return
				}
			*/

			bearerToken := r.Header.Get("Authorization")
			token := strings.Split(bearerToken, " ")
			if len(token) != 2 {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error.message": "%v"}`, "missing token")
				return
			}
			jwt := service.Jwt{
				TokenLifeTime:        a.Conf().Server.TokenLifeTime,
				TokenRefreshLifeTime: a.Conf().Server.TokenRefreshLifeTime,
				TokenRefreshAllowd:   a.Conf().Server.TokenRefreshAllowed,
				TokenKey:             a.Conf().Server.TokenKey,
			}

			user, err := jwt.ValidateToken(token[1])
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error.message": "%v"}`, "invalid token")
				return
			}
			log.Debug("### USER")
			log.Debug(user)

			a.SetUser(user)
			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}

}
