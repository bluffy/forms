package app

import (
	"encoding/json"
	"net/http"

	"github.com/bluffy/forms/models"
	"github.com/bluffy/forms/repository"
	"github.com/bluffy/forms/server/service"
	"github.com/bluffy/forms/util/tools"
	log "github.com/sirupsen/logrus"
)

// HandlerLogin  godoc
// @Tags         public
// @Description  login
// @Accept       json
// @Produce      json
// @Param data body models.UserLoginForm  true "Email & Password"
// @Success      200 {object} models.Token
// @Failure      422 {object} models.AppError
// @Failure      500 {object} models.AppError
// @Router       /api/v1/login [post]
func (app *App) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	form := &models.UserLoginForm{}
	if app.checkForm(form, w, r) {
		return
	}

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 200, err, "")
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.printError(w, http.StatusUnprocessableEntity, 201, err, "")
		return
	}

	jwt := service.Jwt{}
	userDto := user.ToDto()
	token, err := jwt.CreateToken(*userDto)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 202, err, "")
		return
	}

	//dtos := token.ToDto(user)
	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Warn(err)
		app.printError(w, http.StatusInternalServerError, 102, err, "")
	}

}

func (app *App) RefreshLoginToken(w http.ResponseWriter, r *http.Request) {

	token := models.Token{}

	log.Debug("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 103, err, "")
		return
	}

	jwt := service.Jwt{}
	user, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 203, err, "")
		return
	}

	token, err = jwt.CreateToken(user)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 204, err, "")
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 102, err, "")
		return
	}
}
