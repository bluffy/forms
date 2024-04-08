package app

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

func (app *App) HandlerIndex(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	logrus.Println("TEST")
	//res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	/*
		user, err := repository.GetUserByEmail(app.db, "system@bluffy.de")
		if err != nil {
			printError(app, res, http.StatusInternalServerError, "user & password not matched", err)
			return
		}

		if err := json.NewEncoder(res).Encode(user); err != nil {
			logrus.Warn(err)
			printError(app, res, http.StatusInternalServerError, appErrJsonCreationFailure, err)
		}
	*/
}
