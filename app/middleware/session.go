package middleware

import (
	"context"
	"goapp/app"
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func SetLocale(a *app.App) func(next http.Handler) http.Handler {
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

			localizer := i18n.NewLocalizer(a.GetBundle(), locale, accept)
			ctx := context.WithValue(r.Context(), app.ContextLocalizerKey{}, localizer)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}

}

func SetSession(a *app.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			sessStore := session.GetSession(r)
			ctx := context.WithValue(r.Context(), app.ContextUserIDKey{}, &sessStore)
			next.ServeHTTP(w, r.WithContext(ctx))

			/*
				ctx := context.WithValue(ctx, app.ContextSessionStoreKey{}, &sessStore)


				user_id := sessStore.Get("user_id")
				if user_id != nil {
					ctx2 := context.WithValue(ctx1, app.ContextUserIDKey{}, user_id)
					next.ServeHTTP(w, r.WithContext(ctx2))

					//next.ServeHTTP(w, r.WithContext(ctx))
					//return
				} else {
					next.ServeHTTP(w, r.WithContext(ctx1))

				}
			*/

		}
		return http.HandlerFunc(fn)
	}

}
func CheckUserLogin(a *app.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			user_id := r.Context().Value(app.ContextUserIDKey{})
			if user_id == nil || user_id == "" {
				//a.JsonError(w, http.StatusUnauthorized, a.GetLocale("").Text.Session__error_sessen_not_exists_or_expired, nil, false, "")
				a.JsonError(w, http.StatusUnauthorized, "test", nil, false, "")
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

}
