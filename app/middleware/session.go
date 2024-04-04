package middleware

import (
	"context"
	"fmt"
	"net/http"

	"gitea.com/go-chi/session"
)

type UserIDKey struct{}

func SessionCheck() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			sessStore := session.GetSession(r)
			user_id := sessStore.Get("user_id")
			if user_id == nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, `{"error": {"message": "%v"}}`, "session not exists or expired")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey{}, user_id)

			next.ServeHTTP(w, r.WithContext(ctx))

		}
		return http.HandlerFunc(fn)
	}

}
