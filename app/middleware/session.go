package middleware

import (
	"context"
	"goapp/app"
	"log"
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func SetSession(a *app.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			accept := r.Header.Get("Accept-Language")

			locale := r.URL.Query().Get("bl-locale")

			if locale == "" {
				cookie, err := r.Cookie("bl-locale")
				if err == nil {
					locale = cookie.Value
				}
			}

			log.Println(locale)
			localizer := i18n.NewLocalizer(a.GetBundle(), locale, accept)
			ctx := context.WithValue(r.Context(), app.LocalizerKey{}, localizer)

			sessStore := session.GetSession(r)
			user_id := sessStore.Get("user_id")
			if user_id != nil {
				ctx = context.WithValue(ctx, app.UserIDKey{}, user_id)
				//next.ServeHTTP(w, r.WithContext(ctx))
				//return
			}

			next.ServeHTTP(w, r.WithContext(ctx))

		}
		return http.HandlerFunc(fn)
	}

}
func CheckUserLogin(a *app.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user_id := r.Context().Value(app.UserIDKey{})
			if user_id == nil || user_id == "" {
				a.JsonError(w, http.StatusUnauthorized, a.GetLocale("").Text.Session__error_sessen_not_exists_or_expired, nil, false, "")
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

}
