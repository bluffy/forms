package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"goapp/models"
	"goapp/repository"
	"goapp/util/tools"

	"gitea.com/go-chi/session"
	"github.com/go-chi/chi"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
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
// @Router       /bl-api/page/v1/user/login [post]
func (app *App) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	localizer := GetLocalizer(r)

	form := &models.UserLoginForm{}
	if app.checkForm(localizer, form, w, r, nil) {
		return
	}

	msgUserPasswordNotMatched, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerLogin.Error.UserNotExits",
			Other: "user & password not matched",
		},
	})

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		app.ServerLogByRequest(r, err, msgUserPasswordNotMatched, false, "")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserPasswordNotMatched)
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.ServerLogByRequest(r, err, msgUserPasswordNotMatched, false, "!tools.CheckPasswordHash(form.Password, user.Password)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserPasswordNotMatched)
		return
	}

	sess := session.GetSession(r)
	sess.Set("user_id", user.ID)

	_, err = session.RegenerateSession(w, r)
	if err != nil {
		app.ServerLogByRequest(r, err, msgUserPasswordNotMatched, true, "session.RegenerateSession(w, r)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
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
// @Router       /bl-api/page/v1/user/register [post]
func (app *App) HandlerRegister(w http.ResponseWriter, r *http.Request) {

	localizer := GetLocalizer(r)
	form := &models.RegisterUserForm{}
	if app.checkForm(localizer, form, w, r, nil) {
		return
	}

	msgUserExists, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerRegister.Error.UserAlreadyRegistered",
			Other: "user alread exists",
		},
	})

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.ServerLogByRequest(r, err, msgUserExists, true, "")
			app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserExists)
			return
		}
	} else {
		if user.ID != "" {
			app.ServerLogByRequest(r, nil, msgUserExists, true, "user.ID != \"\"")
			app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserExists)
			return
		}

	}
	registerUser, err := form.ToModel()
	if err != nil {
		app.ServerLogByRequest(r, err, "", true, "form.ToModel()")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}
	registerUser, err = repository.CreateRegisterUser(app.db, registerUser)
	if err != nil {
		app.ServerLogByRequest(r, err, "", true, "repository.CreateRegisterUser(app.db, registerUser)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}

	userLink := registerUser.ToLinkModel()
	userLinkBytes, err := json.Marshal(userLink)
	if err != nil {
		app.ServerLogByRequest(r, err, "", true, "json.Marshal(userLink)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}

	userLinkEncrypted, err := tools.EncryptBase64(string(userLinkBytes), app.conf.EncryptKey)
	if err != nil {
		app.ServerLogByRequest(r, err, "", true, "tools.EncryptBase64(string(userLinkBytes)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
	}

	link := app.conf.Server.PublicURL + "/p/user/register/" + userLinkEncrypted
	err = app.sendMail(true, "New Registration", "register link: "+link, "", registerUser.Email)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerRegister.Error.ConfirmationMailNotSend",
				Other: "no confirmation email could be sent",
			},
		})
		app.ServerLogByRequest(r, err, msg, true, "app.sendMail()")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	var response ApiPageResponse
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerRegister.RegisteringSuccessful",
			Other: "register was successful, check your email please!t",
		},
	})
	response.Message = &msg
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.ServerLogByRequest(r, err, "", true, "form.ToModel()")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
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
// @Router       /bl-api/page/v1/user/register [get]
func (app *App) HandlerRegisterLink(w http.ResponseWriter, r *http.Request) {

	link := chi.URLParam(r, "link")
	localizer := GetLocalizer(r)

	if localizer == nil {
		app.ServerLogByRequest(r, nil, DEFAULT_ERROR, true, "")
		app.RenderError(r, w, &PageError{Message: DEFAULT_ERROR})
		return
	}

	errMsgLinkInvalid, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerRegisterLink.Register_link_is_invalid",
			Other: "the link is invalid or expired",
		},
	})

	logrus.Print("link")
	logrus.Print(link)
	linkDec, err := tools.DecryptBase64(link, app.conf.EncryptKey)
	if err != nil {
		app.ServerLogByRequest(r, err, errMsgLinkInvalid, true, "tools.DecryptBase64(link, app.conf.EncryptKey)")
		app.RenderError(r, w, &PageError{Message: errMsgLinkInvalid})
		return
	}

	linkUser := &models.RegisterUserLink{}

	err = json.Unmarshal([]byte(linkDec), linkUser)
	if err != nil {
		app.ServerLogByRequest(r, err, errMsgLinkInvalid, true, "json.Unmarshal([]byte(linkDec), linkUser)")
		app.RenderError(r, w, &PageError{Message: errMsgLinkInvalid})
		return
	}

	dbRegisterUser, err := repository.ReadRegisterUser(app.db, linkUser.ID)
	if err != nil {
		app.ServerLogByRequest(r, err, errMsgLinkInvalid, true, "repository.ReadRegisterUser(app.db, linkUser.ID)")
		app.RenderError(r, w, &PageError{Message: errMsgLinkInvalid})
		return
	}
	if dbRegisterUser.UpdatedAt != linkUser.CreatedAt {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerRegisterLink.Register_link_is_expired",
				Other: "the link is expired",
			},
		})
		app.ServerLogByRequest(r, err, msg, false, "dbRegisterUser.UpdatedAt != linkUser.CreatedAt")
		app.RenderError(r, w, &PageError{Message: msg})
		return
	}

	dbUser, err := repository.GetUserByEmail(app.db, dbRegisterUser.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			msg, _ := localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID: LOCALE_MSG_ID__COMMON_SERVER_ERROR,
				},
			})
			app.ServerLogByRequest(r, err, msg, true, "repository.GetUserByEmail(app.db, dbRegisterUser.Email)")
			app.RenderError(r, w, &PageError{Message: msg})
			return
		}
	}

	if dbUser != nil && dbUser.Email == dbRegisterUser.Email {

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerRegisterLink.User_already_registered",
				Other: "user alread exists",
			},
		})
		app.ServerLogByRequest(r, nil, msg, true, "dbUser != nil && dbUser.Email == dbRegisterUser.Email")
		app.RenderError(r, w, &PageError{Message: msg})

		return
	}

	dbUser = dbRegisterUser.ToUserModel()
	_, err = repository.CreateUser(app.db, dbUser)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: LOCALE_MSG_ID__COMMON_SERVER_ERROR,
			},
		})
		app.ServerLogByRequest(r, err, msg, true, "repository.CreateUser(app.db, dbUser)")
		app.RenderError(r, w, &PageError{Message: msg})

		return
	}
	err = repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)
	if err != nil {
		app.ServerLogByRequest(r, err, "", true, "repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)")
	}

	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerRegisterLink.Message_user_created",
			Other: "successful, you can login now!",
		},
	})
	page := Page{
		Message: &msg,
	}
	app.RenderPage(r, w, &page)

}

/*
func (app *App) RefreshLoginToken(w http.ResponseWriter, r *http.Request) {

	token := models.Token{}

	logrus.Debug("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		app.JsonError(r,w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__json_decode, "Error__json_decode", err)
		return
	}

	jwt := service.Jwt{}
	user, sessionId, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		app.JsonError(r,w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__invalid_token, "Page_auth__error__invalid_token", err)

		return
	}

	token, err = jwt.CreateToken(*user, *sessionId)
	if err != nil {
		app.JsonError(r,w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Page_auth__error__unable_create_access_token, "Page_auth__error__unable_create_access_token", err)
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		app.JsonError(r,w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__json_encode, "Error__json_encode", err)
		return
	}
}
*/

// HandlerLogin  godoc
// @Tags         public
// @Description  index test
// @Accept       json
// @Produce      json
// @Param data body models.UserLoginForm  true "Email & Password"
// @Success      204 {object} models.Token
// @Failure      401 {object} models.AppError
// @Failure      422 {object} models.AppError
// @Failure      500 {object} models.AppError "Response JSON"
// @Router       /bl-api/page/v1/user/forgot_password[post]
func (app *App) HandlerGeneratePasswordLink(w http.ResponseWriter, r *http.Request) {

}
