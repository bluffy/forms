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

			/*
				bearerToken := r.Header.Get("Authorization")
				token := strings.Split(bearerToken, " ")
				if len(token) != 2 {
					w.WriteHeader(http.StatusUnauthorized)
					fmt.Fprintf(w, `{"error": {"message": "%v"}}`, "missing token")
					return
				}
			*/
			token := TokenFromCookie(r)
			if token == "" {
				token = TokenFromHeader(r)
			}
			if token == "" {
				token = TokenFromQuery(r)
			}

			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error": {"message": "%v"}}`, "missing token")
				return
			}

			jwt := service.Jwt{
				TokenLifeTime:        a.Conf().Server.TokenLifeTime,
				TokenRefreshLifeTime: a.Conf().Server.TokenRefreshLifeTime,
				TokenRefreshAllowd:   a.Conf().Server.TokenRefreshAllowed,
				TokenKey:             a.Conf().Server.TokenKey,
			}

			user, _, err := jwt.ValidateToken(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error": {"message": "%v"}}`, "invalid token")
				return
			}
			log.Debug("### USER")
			log.Debug(user)

			a.SetUser(*user)
			next.ServeHTTP(w, r)

		}
		return http.HandlerFunc(fn)
	}

}

// TokenFromCookie tries to retreive the token string from a cookie named
// "jwt".
func TokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// TokenFromHeader tries to retreive the token string from the
// "Authorization" reqeust header: "Authorization: BEARER T".
func TokenFromHeader(r *http.Request) string {
	// Get token from authorization header.
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}

// TokenFromQuery tries to retreive the token string from the "jwt" URI
// query parameter.
//
// To use it, build our own middleware handler, such as:
//
//	func Verifier(ja *JWTAuth) func(http.Handler) http.Handler {
//		return func(next http.Handler) http.Handler {
//			return Verify(ja, TokenFromQuery, TokenFromHeader, TokenFromCookie)(next)
//		}
//	}
func TokenFromQuery(r *http.Request) string {
	// Get token from query param named "jwt".
	return r.URL.Query().Get("jwt")
}
