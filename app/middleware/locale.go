package middleware

import (
	"context"
	"goapp/app"
	"net/http"

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
