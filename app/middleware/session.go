package middleware

import (
	"context"
	"goapp/app"
	"net/http"

	"gitea.com/go-chi/session"
)

type UserIDKey struct{}

func SessionCheck(a *app.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			sessStore := session.GetSession(r)
			user_id := sessStore.Get("user_id")
			if user_id == nil {
				a.JsonError(w, http.StatusUnauthorized, a.GetLocale("").Text.Session__error_sessen_not_exists_or_expired, nil, false, "")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey{}, user_id)

			next.ServeHTTP(w, r.WithContext(ctx))

		}
		return http.HandlerFunc(fn)
	}

}
