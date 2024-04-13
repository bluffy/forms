package middleware

import (
	"goapp/app"
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

/*
	func SetSession(a *app.App) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			fn := func(w http.ResponseWriter, r *http.Request) {

				sessStore := session.GetSession(r)
				user_id := sessStore.Get(app.SessionKeyUserID)
				if user_id != nil {
					ctx = context.WithValue(ctx, app.ContextUserIDKey{}, user_id)

				}
				next.ServeHTTP(w, r.WithContext(ctx))

			}
			return http.HandlerFunc(fn)
		}

}
*/
func CheckUserLogin(a *app.App) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			store := session.GetSession(r)
			user_id := store.Get(app.SessionKeyUserID)
			//user_id := r.Context().Value(app.ContextUserIDKey{})
			if user_id == nil || user_id == "" {
				localizer := app.GetLocalizer(r)
				msg := "session not exists or expired!"
				if localizer == nil {
					msg, _ = app.GetLocalizer(r).Localize(&i18n.LocalizeConfig{
						DefaultMessage: &i18n.Message{
							ID:    "Api.SessionCheck.sessen_not_exists_or_expire",
							Other: "session not exists or expired!",
						},
					})
				}
				a.ErrorRequestLog(r, nil, msg, false, "")

				a.JsonError(r, w, http.StatusUnauthorized, msg)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

}
