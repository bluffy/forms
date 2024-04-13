package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"goapp/models"
	"goapp/repository"
	"goapp/util/tools"

	"gitea.com/go-chi/session"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// HandleAuthLoginForm  godoc
// @Tags         auth
// @Description  login
// @Accept       json
// @Produce      json
// @Param data body models.UserLoginForm  true "Email & Password"
// @Success      204
// @Failure      422 {object} app.ErrResponse
// @Router       /bl-api/v1/user/login [post]
func (app *App) HandleAuthLoginForm(w http.ResponseWriter, r *http.Request) {
	localizer := GetLocalizer(r)

	form := &models.UserLoginForm{}
	weiter, err := app.checkForm(localizer, form, w, r, nil)
	if err != nil {
		msg := GetMessageServerError(r, "(forms)")
		app.ErrorRequestLog(r, err, msg, true, msg)
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	} else if !weiter {
		return
	}

	msgUserPasswordNotMatched, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandleAuthLoginForm.Error.UserNotExits",
			Other: "user & password not matched",
		},
	})

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		app.ErrorRequestLog(r, err, msgUserPasswordNotMatched, false, "")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserPasswordNotMatched)
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.ErrorRequestLog(r, err, msgUserPasswordNotMatched, false, "!tools.CheckPasswordHash(form.Password, user.Password)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserPasswordNotMatched)
		return
	}

	session.GetSession(r).Set(SessionKeyUserID, user.ID)

	_, err = session.RegenerateSession(w, r)
	if err != nil {
		app.ErrorRequestLog(r, err, msgUserPasswordNotMatched, true, "session.RegenerateSession(w, r)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

// HandleAuthRegisterFrom  godoc
// @Tags         auth
// @Description  login
// @Accept       json
// @Produce      json
// @Param data body models.RegisterUserForm true "register informations"
// @Success      200 {object} app.ApiPageResponse
// @Failure      422 {object} app.ErrResponse
// @Failure      500 {object} app.ErrResponse "Response JSON"
// @Router       /bl-api/v1/user/register [post]
func (app *App) HandleAuthRegisterFrom(w http.ResponseWriter, r *http.Request) {

	localizer := GetLocalizer(r)
	form := &models.RegisterUserForm{}

	weiter, err := app.checkForm(localizer, form, w, r, nil)

	if err != nil {
		msg := GetMessageServerError(r, "(forms)")
		app.ErrorRequestLog(r, err, msg, true, msg)
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
	} else if !weiter {
		return
	}

	msgUserExists, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandleAuthRegisterFrom.Error.UserAlreadyRegistered",
			Other: "user alread exists",
		},
	})

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.ErrorRequestLog(r, err, msgUserExists, true, "")
			app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserExists)
			return
		}
	} else {
		if user.ID != "" {
			app.ErrorRequestLog(r, nil, msgUserExists, true, "user.ID != \"\"")
			app.JsonError(r, w, http.StatusUnprocessableEntity, msgUserExists)
			return
		}

	}
	registerUser, err := form.ToModel()
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "form.ToModel()")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}
	registerUser, err = repository.CreateRegisterUser(app.db, registerUser)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "repository.CreateRegisterUser(app.db, registerUser)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}

	userLink := registerUser.ToLinkModel()
	userLinkBytes, err := json.Marshal(userLink)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.Marshal(userLink)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}

	userLinkEncrypted, err := tools.EncryptBase64(string(userLinkBytes), app.conf.EncryptKey)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "tools.EncryptBase64(string(userLinkBytes)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
	}
	link := app.conf.Server.ClientUrl + "/user/register/" + userLinkEncrypted

	msgMail, _ := localizer.Localize(&i18n.LocalizeConfig{
		TemplateData: map[string]interface{}{
			"Name":         registerUser.FirstName + " " + registerUser.LastName,
			"RegisterLink": link,
		},
		DefaultMessage: &i18n.Message{
			ID: "Api.HandleAuthRegisterFrom.Mail.RegistrationLink",
			Other: `Hello {{.Name}},
Please click the following link within 3 hours to register your account:

{{.RegisterLink}}

link is not working? Try copying and pasting it into the browser.`,
		},
	})

	logrus.Debug(link)
	err = app.sendMail(true, "New Registration", msgMail, "", registerUser.Email)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandleAuthRegisterFrom.Error.ConfirmationMailNotSend",
				Other: "no confirmation email could be sent",
			},
		})
		app.ErrorRequestLog(r, err, msg, true, "app.sendMail()")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	var response ApiPageResponse
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		TemplateData: map[string]interface{}{
			"UserEmail": form.Email,
		},
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandleAuthRegisterFrom.RegisteringSuccessful",
			Other: "A confirmation email was sent to **{{.UserEmail}}**. Please check your mailbox within 3 hours to complete the registration process.",
		},
	})
	response.Message = &msg
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.NewEncoder(w).Encode(response);")
		app.JsonError(r, w, http.StatusUnprocessableEntity, GetMessageServerError(r, "(json encode)"))
	}

}

// HandleAuthLogout  godoc
// @Tags         auth
// @Description  login
// @Accept       json
// @Success      204
// @Failure      422 {object} app.ErrResponse
// @Failure      500 {object} app.ErrResponse "Response JSON"
// @Router       /bl-api/v1/user/logout [get]
func (app *App) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {

	//localizer := GetLocalizer(r)

	sess := session.GetSession(r)
	userId := sess.Get(SessionKeyUserID).(string)
	msg := GetMessageServerError(r, "(database)")

	userDb, err := repository.ReadUser(app.db, userId)
	if err != nil {
		app.ErrorRequestLog(r, err, msg, true, "repository.ReadUser(app.db, userId)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
	}
	_, err = repository.UpdateUser(app.db, userDb)
	if err != nil {
		app.ErrorRequestLog(r, err, msg, true, "repository.UpdateUser(app.db, userDb)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
	}

	err = session.GetSession(r).Destroy(w, r)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.NewEncoder(w).Encode(response);")
		app.JsonError(r, w, http.StatusUnprocessableEntity, GetMessageServerError(r, "(json encode)"))
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

/*

func (app *App) HandlerRegisterLinkGet(w http.ResponseWriter, r *http.Request) {

	link := chi.URLParam(r, "link")
	localizer := GetLocalizer(r)

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
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "tools.DecryptBase64(link, app.conf.EncryptKey)")
		app.RenderError(r, w, &PageError{Message: errMsgLinkInvalid})
		return
	}

	linkUser := &models.RegisterUserLink{}

	err = json.Unmarshal([]byte(linkDec), linkUser)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "json.Unmarshal([]byte(linkDec), linkUser)")
		app.RenderError(r, w, &PageError{Message: errMsgLinkInvalid})
		return
	}

	dbRegisterUser, err := repository.ReadRegisterUser(app.db, linkUser.ID)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "repository.ReadRegisterUser(app.db, linkUser.ID)")
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
		app.ErrorRequestLog(r, err, msg, false, "dbRegisterUser.UpdatedAt != linkUser.CreatedAt")
		app.RenderError(r, w, &PageError{Message: msg})
		return
	}

	dbUser, err := repository.GetUserByEmail(app.db, dbRegisterUser.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			msg := GetMessageServerError(r, "(database)")
			app.ErrorRequestLog(r, err, msg, true, "repository.GetUserByEmail(app.db, dbRegisterUser.Email)")
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
		app.ErrorRequestLog(r, nil, msg, true, "dbUser != nil && dbUser.Email == dbRegisterUser.Email")
		app.RenderError(r, w, &PageError{Message: msg})

		return
	}

	dbUser = dbRegisterUser.ToUserModel()
	_, err = repository.CreateUser(app.db, dbUser)
	if err != nil {
		msg := GetMessageServerError(r, "(database)")
		app.ErrorRequestLog(r, err, msg, true, "repository.CreateUser(app.db, dbUser)")
		app.RenderError(r, w, &PageError{Message: msg})

		return
	}
	err = repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)")
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
*/

func (app *App) HanderCreateUserFromMailLink(w http.ResponseWriter, r *http.Request) {

	localizer := GetLocalizer(r)

	var form struct {
		Link string `json:"decoded" form:"required"`
	}

	errMsgLinkInvalid, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HanderCreateUserFromMailLink.Register_link_is_invalid",
			Other: "the link is invalid or expired",
		},
	})

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil || form.Link == "" {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "json.NewDecoder(r.Body).Decode(form)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}

	linkDec, err := tools.DecryptBase64(form.Link, app.conf.EncryptKey)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "tools.DecryptBase64(link, app.conf.EncryptKey)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}

	linkUser := &models.RegisterUserLink{}

	err = json.Unmarshal([]byte(linkDec), linkUser)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "json.Unmarshal([]byte(linkDec), linkUser)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}

	dbRegisterUser, err := repository.ReadRegisterUser(app.db, linkUser.ID)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "repository.ReadRegisterUser(app.db, linkUser.ID)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}
	if dbRegisterUser.UpdatedAt != linkUser.CreatedAt {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HanderCreateUserFromMailLink.Register_link_is_expired",
				Other: "the link is expired",
			},
		})
		app.ErrorRequestLog(r, err, msg, false, "dbRegisterUser.UpdatedAt != linkUser.CreatedAt")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	dbUser, err := repository.GetUserByEmail(app.db, dbRegisterUser.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			msg := GetMessageServerError(r, "(database)")
			app.ErrorRequestLog(r, err, msg, true, "repository.GetUserByEmail(app.db, dbRegisterUser.Email)")
			app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
			return
		}
	}

	if dbUser != nil && dbUser.Email == dbRegisterUser.Email {

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HanderCreateUserFromMailLink.User_already_registered",
				Other: "user alread exists",
			},
		})
		app.ErrorRequestLog(r, nil, msg, true, "dbUser != nil && dbUser.Email == dbRegisterUser.Email")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)

		return
	}

	dbUser = dbRegisterUser.ToUserModel()
	_, err = repository.CreateUser(app.db, dbUser)
	if err != nil {
		msg := GetMessageServerError(r, "(database)")
		app.ErrorRequestLog(r, err, msg, true, "repository.CreateUser(app.db, dbUser)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)

		return
	}
	err = repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)")
	}

	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HanderCreateUserFromMailLink.Message_user_created",
			Other: "successful, you can login now!",
		},
	})

	var response ApiPageResponse

	response.Message = &msg
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.NewEncoder(w).Encode(response);")
		app.JsonError(r, w, http.StatusUnprocessableEntity, GetMessageServerError(r, "(json encode)"))
	}

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

func (app *App) HandlerGenerateMailWithPasswordLink(w http.ResponseWriter, r *http.Request) {

	localizer := GetLocalizer(r)

	type UserPasswordLinkForm struct {
		Email string `json:"email" form:"required,max=255,email"`
	}

	form := &UserPasswordLinkForm{}
	weiter, err := app.checkForm(localizer, form, w, r, nil)

	if err != nil {
		msg := GetMessageServerError(r, "(forms)")
		app.ErrorRequestLog(r, err, msg, true, msg)
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
	} else if !weiter {
		return
	}

	user, err := repository.GetUserByEmail(app.db, form.Email)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerGenerateMailWithPasswordLink.Error.UserNotExits",
				Other: "user not exists",
			},
		})
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			app.ErrorRequestLog(r, err, msg, false, "")
		}
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	now := time.Now()
	user.NewPasswordRequest = &now
	user, err = repository.UpdateUser(app.db, user)
	if err != nil {
		msg := GetMessageServerError(r, "(database)")
		app.ErrorRequestLog(r, err, msg, true, msg)
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	passwordLink := &models.UserPasswordLink{
		ID:                 user.ID,
		NewPasswordRequest: *user.NewPasswordRequest,
	}

	passwordLinkBytes, err := json.Marshal(passwordLink)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.Marshal(userLink)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
		return
	}

	userLinkEncrypted, err := tools.EncryptBase64(string(passwordLinkBytes), app.conf.EncryptKey)
	if err != nil {
		app.ErrorRequestLog(r, err, "", true, "tools.EncryptBase64(string(userLinkBytes)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, "")
	}

	msgTitle, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerGenerateMailWithPasswordLink.Mail.Subject",
			Other: "recover your account",
		},
	})
	msgText, _ := localizer.Localize(&i18n.LocalizeConfig{
		TemplateData: map[string]interface{}{
			"Name":        strings.TrimSpace(user.FirstName + " " + user.LastName),
			"RecoverLink": app.conf.Server.ClientUrl + "/user/forgot_password/" + userLinkEncrypted,
		},
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerGenerateMailWithPasswordLink.Mail.Text",
			Other: "Hello {{.Name}},\n\nPlease click the following link within 3 hours to restore your account:\n\n{{.RecoverLink}}\n\nlink is not working? Try copying and pasting it into the browser.",
		},
	})

	err = app.sendMail(true, msgTitle, msgText, "", user.Email)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerGenerateMailWithPasswordLink.Error.ConfirmationMailNotSend",
				Other: "no confirmation email could be sent",
			},
		})
		app.ErrorRequestLog(r, err, msg, true, "app.sendMail()")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	var response ApiPageResponse
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		TemplateData: map[string]interface{}{
			"UserEmail": user.Email,
		},
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerGenerateMailWithPasswordLink.Successful",
			Other: "A confirmation email was sent to **{{.UserEmail}}**. Please check your mailbox within 3 hours to complete the restore process.",
		},
	})
	response.Message = &msg
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.NewEncoder(w).Encode(response);")
		app.JsonError(r, w, http.StatusUnprocessableEntity, GetMessageServerError(r, "(json encode)"))
	}
}

func (app *App) HanldeCheckUser(res http.ResponseWriter, req *http.Request) {

	res.WriteHeader(http.StatusNoContent)
	res.Write([]byte(""))

}

func (app *App) HandlerCreateNewPasswordFromMailLink(w http.ResponseWriter, r *http.Request) {

	localizer := GetLocalizer(r)

	form := &models.UserPasswordForm{}

	weiter, err := app.checkForm(localizer, form, w, r, nil)

	if err != nil {
		msg := GetMessageServerError(r, "(forms)")
		app.ErrorRequestLog(r, err, msg, true, msg)
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	} else if !weiter {
		return
	}

	errMsgLinkInvalid, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerCreateNewPasswordFromMailLink.Link_is_invalid",
			Other: "the link is invalid or expired",
		},
	})

	if form.Link == "" {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "form.Link == ''")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}

	linkDec, err := tools.DecryptBase64(form.Link, app.conf.EncryptKey)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "tools.DecryptBase64(link, app.conf.EncryptKey)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}

	linkUser := &models.UserPasswordLink{}

	err = json.Unmarshal([]byte(linkDec), linkUser)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "json.Unmarshal([]byte(linkDec), linkUser)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}

	dbUser, err := repository.ReadUser(app.db, linkUser.ID)
	if err != nil {
		app.ErrorRequestLog(r, err, errMsgLinkInvalid, true, "repository.ReadRegisterUser(app.db, linkUser.ID)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, errMsgLinkInvalid)
		return
	}
	if *dbUser.NewPasswordRequest != linkUser.NewPasswordRequest {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerCreateNewPasswordFromMailLink.Link_expired",
				Other: "the link is expired",
			},
		})
		app.ErrorRequestLog(r, err, msg, true, "dbUser.NewPasswordRequest != &linkUser.NewPasswordRequest")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
		return
	}

	dbUser.Password, err = tools.HashPassword(form.Password)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerCreateNewPasswordFromMailLink.Password_is_malformed",
				Other: "selected password is malformed",
			},
		})
		app.ErrorRequestLog(r, err, msg, true, "tools.HashPassword(form.Password)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
	}

	now := time.Now()
	dbUser.NewPasswordRequest = &now
	_, err = repository.UpdateUser(app.db, dbUser)
	if err != nil {
		msg := GetMessageServerError(r, "(database)")
		app.ErrorRequestLog(r, err, msg, true, "UpdateUser(app.db, dbUser)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msg)
	}

	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerCreateNewPasswordFromMailLink.Password_successful_changed",
			Other: "successful, you can login now!",
		},
	})

	var response ApiPageResponse

	response.Message = &msg
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.ErrorRequestLog(r, err, "", true, "json.NewEncoder(w).Encode(response);")
		app.JsonError(r, w, http.StatusUnprocessableEntity, GetMessageServerError(r, "(json encode)"))
	}

}
