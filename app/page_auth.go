package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"goapp/models"
	"goapp/repository"
	"goapp/util/tools"

	"gitea.com/go-chi/session"
	"gorm.io/gorm"
)

// HandlerLogin  godoc
// @Tags         public
// @Description  login
// @Accept       json
// @Produce      json
// @Param data body models.UserLoginForm  true "Email & Password"
// @Success      204 {object} models.Token
// @Failure      401 {object} models.AppError
// @Failure      422 {object} models.AppError
// @Failure      500 {object} models.AppError "Response JSON"
// @Router       /bl-api/page/v1/login [post]
func (app *App) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	form := &models.UserLoginForm{}
	if app.checkForm("", form, w, r, nil) {
		return
	}

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__user_not_exists, "Page_auth__error__user_not_exists", err)
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__wrong_password, "Page_auth__error__wrong_password", err)
		return
	}

	sess := session.GetSession(r)
	sess.Set("user_id", user.ID)

	_, err = session.RegenerateSession(w, r)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__session_regeneration, "Page_auth__error__session_regeneration", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

// HandlerLogin  godoc
// @Tags         public
// @Description  login
// @Accept       json
// @Produce      json
// @Param data body models.RegisterUserForm  true "Email & Password"
// @Success      200 {object} models.RegisterUserForm
// @Failure      401 {object} models.AppError
// @Failure      422 {object} models.AppError
// @Failure      500 {object} models.AppError "Response JSON"
// @Router       /bl-api/page/v1/register [post]
func (app *App) HandlerRgister(w http.ResponseWriter, r *http.Request) {

	form := &models.RegisterUserForm{}
	if app.checkForm("", form, w, r, nil) {
		return
	}
	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__user_already_registered, "Page_auth__error__user_already_registered", err)
			return
		}
	} else {
		if user.ID != "" {
			app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__user_already_registered, "Page_auth__error__user_already_registered", err)
			return
		}

	}
	registerUser, err := form.ToModel()
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__commen_server_error, "Error__commen_server_error", err)
		return
	}
	registerUser, err = repository.CreateRegisterUser(app.db, registerUser)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__database_error, "Error__database_error", err)
		return
	}

	mail_text := "register link: " + registerUser.ID
	var mail models.Mail
	mail.Status = 0
	mail.Sender = "system@bluffy.de"
	mail.Text = &mail_text
	mail.Recipient = registerUser.Email
	_, err = repository.CreateMail(app.db, &mail)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__database_error, "Error__database_error in CreateMail", err)
		return
	}

	// SEND MAIL

	var response models.UserRegisterResponse
	response.Message = app.GetLocale("").Text.Page_auth__regsitering_success
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__database_error, "Error__json_encode", err)
	}

}

/*
func (app *App) RefreshLoginToken(w http.ResponseWriter, r *http.Request) {

	token := models.Token{}

	log.Debug("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		app.jsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__json_decode, "Error__json_decode", err)
		return
	}

	jwt := service.Jwt{}
	user, sessionId, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		app.jsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__invalid_token, "Page_auth__error__invalid_token", err)

		return
	}

	token, err = jwt.CreateToken(*user, *sessionId)
	if err != nil {
		app.jsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__unable_create_access_token, "Page_auth__error__unable_create_access_token", err)
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		app.jsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__json_encode, "Error__json_encode", err)
		return
	}
}
*/
