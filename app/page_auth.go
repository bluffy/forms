package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"goapp/models"
	"goapp/repository"
	"goapp/util/tools"

	"gitea.com/go-chi/session"
	"github.com/go-chi/chi/v5"
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
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__user_not_exists, err, false, "")
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__wrong_password, err, false, "")
		return
	}

	sess := session.GetSession(r)
	sess.Set("user_id", user.ID)

	_, err = session.RegenerateSession(w, r)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__session_regeneration, err, false, "")
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
			app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__user_already_registered, err, false, "")
			return
		}
	} else {
		if user.ID != "" {
			app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__user_already_registered, err, false, "")
			return
		}

	}
	registerUser, err := form.ToModel()
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__commen_server_error, err, false, "")
		return
	}
	registerUser, err = repository.CreateRegisterUser(app.db, registerUser)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__database_error, err, false, "")
		return
	}

	userLink := registerUser.ToLinkModel()

	userLinkBytes, err := json.Marshal(userLink)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__json_encode, err, false, "")
		return
	}

	userLinkEncrypted, err := tools.EncryptBase64(string(userLinkBytes), app.conf.EncryptKey)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, "encyrption error", err, false, "")
	}

	link := app.conf.Server.PublicURL + "/p/register/" + userLinkEncrypted
	err = app.sendMail(true, "New Registration", "register link: "+link, "", registerUser.Email)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, "Register Mail not Sent", err, false, "")
		return
	}

	var response PageResponse
	response.Message = &app.GetLocale("").Text.Page_auth__regsitering_success
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__json_encode, err, false, "")
	}

}

// HandlerLogin  godoc
// @Tags         public
// @Description  login
// @Accept       json
// @Produce      json
// @Param        link  path   string  true  "link encoded"
// @Success      200 {object} models.RegisterUserForm
// @Failure      401 {object} models.AppError
// @Failure      422 {object} models.AppError
// @Failure      500 {object} models.AppError "Response JSON"
// @Router       /bl-api/page/v1/register [get]
func (app *App) HandlerRgisterLink(w http.ResponseWriter, r *http.Request) {

	link := chi.URLParam(r, "link")

	if link == "" {
		w.Write([]byte(app.errMessage("param link is missing", nil, false, "")))
		return
	}

	linkDec, err := tools.DecryptBase64(link, app.conf.EncryptKey)
	if err != nil {
		w.Write([]byte(app.errMessage(app.GetLocale("").Text.Page.Auth.Register_link_is_invalid, nil, false, "")))
		return
	}

	linkUser := &models.RegisterUserLink{}

	err = json.Unmarshal([]byte(linkDec), linkUser)
	if err != nil {
		w.Write([]byte(app.errMessage(app.GetLocale("").Text.Page.Auth.Register_link_is_invalid, nil, false, "")))
		return
	}

	dbRegisterUser, err := repository.ReadRegisterUser(app.db, linkUser.ID)
	if err != nil {
		w.Write([]byte(app.errMessage(app.GetLocale("").Text.Page.Auth.Register_link_is_invalid, nil, false, "")))
		return
	}
	if dbRegisterUser.UpdatedAt != linkUser.CreatedAt {
		w.Write([]byte(app.errMessage(app.GetLocale("").Text.Page.Auth.Register_link_is_expired, nil, false, "")))

		return
	}

	dbUser, err := repository.GetUserByEmail(app.db, dbRegisterUser.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			w.Write([]byte(app.errMessage(app.GetLocale("").Text.Error__database_error, nil, false, "")))

			return
		}
	}

	if dbUser != nil && dbUser.Email == dbRegisterUser.Email {
		w.Write([]byte(app.errMessage(app.GetLocale("").Text.Page.Auth.Register_user_already_registered, nil, false, "")))

		return
	}

	dbUser = dbRegisterUser.ToUserModel()
	_, err = repository.CreateUser(app.db, dbUser)
	if err != nil {
		w.Write([]byte(app.errMessage(app.GetLocale("").Text.Error__database_error, nil, false, "")))

		return
	}
	_ = repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)

	w.Write([]byte(app.errMessage(app.GetLocale("").Text.Page.Auth.Message_user_created, nil, false, "")))

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
