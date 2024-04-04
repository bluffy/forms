package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"goapp/models"
	"goapp/repository"
	"goapp/service"
	"goapp/util/tools"

	"gitea.com/go-chi/session"
	log "github.com/sirupsen/logrus"
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

	sess := session.GetSession(r)
	sess.Set("user_id", user.ID)

	_, err = session.RegenerateSession(w, r)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 202, err, "")
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
	if app.checkForm(form, w, r) {
		return
	}
	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {

			app.printError(w, http.StatusUnprocessableEntity, 210, err, "")
			return
		}
	} else {
		if user.ID != "" {
			app.printError(w, http.StatusUnprocessableEntity, 210, err, "")
			return
		}

	}
	registerUser, err := form.ToModel()
	if err != nil {
		app.printError(w, http.StatusInternalServerError, 105, err, "")
		return
	}
	registerUser, err = repository.CreateRegisterUser(app.db, registerUser)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 104, err, "")
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
		app.printError(w, http.StatusUnprocessableEntity, 104, err, "")
		return
	}

	// SEND MAIL

	var response models.UserRegisterResponse
	response.Message = app.GetLangText("page_register__message_success", "")
	if err := json.NewEncoder(w).Encode(response); err != nil {

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
	user, sessionId, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		app.printError(w, http.StatusUnprocessableEntity, 203, err, "")
		return
	}

	token, err = jwt.CreateToken(*user, *sessionId)
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
