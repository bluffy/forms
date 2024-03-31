package app

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/bluffy/forms/repository"
)

// HandleQueryGet godoc
// @Tags         public
// @Description  get data
// @Accept       json
// @Produce      json
// @Router       /api/v1/intern [get]
// @Security     Token
// @Failure      401 {object} models.AppError
func (app *App) HandlerIntern(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	log.Println("TEST")
	//res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	user, err := repository.GetUserByEmail(app.db, "system@bluffy.de")
	if err != nil {
		printError(app, res, http.StatusInternalServerError, "user & password not matched", err)
		return
	}

	if err := json.NewEncoder(res).Encode(user); err != nil {
		log.Warn(err)
		printError(app, res, http.StatusInternalServerError, appErrJsonCreationFailure, err)
	}

}
