package app

import (
	"net/http"
)

func (app *App) HandlerIntern(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")
	//mysess := session.GetSession(req)
	//userId := mysess.Get("user_id")

	/*
		user, err := repository.ReadUser(app.db, userId.(string))
		if err != nil {
			app.printError(res, http.StatusUnprocessableEntity, 200, err, "")
			return
		}

		if err := json.NewEncoder(res).Encode(user.ToDto()); err != nil {
			logrus.Warn(err)
			app.printError(res, http.StatusInternalServerError, 102, err, "")
		}#
	*/

	//	printError(app, res, http.StatusInternalServerError, "user & password not matched", err)

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
