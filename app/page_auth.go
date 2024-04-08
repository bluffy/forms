package app

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"

	"goapp/models"
	"goapp/repository"
	"goapp/util/tools"

	"gitea.com/go-chi/session"
	"github.com/go-chi/chi/v5"
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
// @Router       /bl-api/page/v1/login [post]
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
		app.JsonError(w, http.StatusUnprocessableEntity, msgUserPasswordNotMatched, err, false, "")
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		app.JsonError(w, http.StatusUnprocessableEntity, msgUserPasswordNotMatched, err, false, "")
		return
	}

	sess := session.GetSession(r)
	sess.Set("user_id", user.ID)

	_, err = session.RegenerateSession(w, r)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, GetDefaultMessage(localizer, DEFAULT_MESSAGE_COMMON_SERVER_ERROR), err, true, "session.RegenerateSession")
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
			app.JsonError(w, http.StatusUnprocessableEntity, msgUserExists, err, false, "")
			return
		}
	} else {
		if user.ID != "" {
			app.JsonError(w, http.StatusUnprocessableEntity, msgUserExists, err, false, "")
			return
		}

	}
	registerUser, err := form.ToModel()
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, GetDefaultMessage(localizer, DEFAULT_MESSAGE_COMMON_SERVER_ERROR), err, true, "form.ToModel()")
		return
	}
	registerUser, err = repository.CreateRegisterUser(app.db, registerUser)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, GetDefaultMessage(localizer, DEFAULT_MESSAGE_COMMON_SERVER_ERROR), err, false, "CreateRegisterUser(app.db, registerUser)")
		return
	}

	userLink := registerUser.ToLinkModel()
	userLinkBytes, err := json.Marshal(userLink)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, GetDefaultMessage(localizer, DEFAULT_MESSAGE_COMMON_SERVER_ERROR), err, true, "json.Marshal(userLink)")
		return
	}

	userLinkEncrypted, err := tools.EncryptBase64(string(userLinkBytes), app.conf.EncryptKey)
	if err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, GetDefaultMessage(localizer, DEFAULT_MESSAGE_COMMON_SERVER_ERROR), err, true, "tools.EncryptBase64(string(userLinkBytes)")
	}

	link := app.conf.Server.PublicURL + "/p/register/" + userLinkEncrypted
	err = app.sendMail(true, "New Registration", "register link: "+link, "", registerUser.Email)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerRegister.Error.ConfirmationMailNotSend",
				Other: "no confirmation email could be sent",
			},
		})
		app.JsonError(w, http.StatusUnprocessableEntity, msg, err, false, "")
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
		app.JsonError(w, http.StatusUnprocessableEntity, GetDefaultMessage(localizer, DEFAULT_MESSAGE_COMMON_SERVER_ERROR), err, false, "")
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
func (app *App) HandlerRegisterLink(w http.ResponseWriter, r *http.Request) {
	link := chi.URLParam(r, "link")
	log.Println("test1", link)

	localizer := GetLocalizer(r)
	//var tmplFile = "templates/page/default.html"
	//tmpl, err := template.New(tmplFile).ParseFS(app.templateFS, tmplFile)
	templateFile := "register-link.html"
	view := template.Must(template.ParseFS(app.templateFS, "templates/page/*"))

	errMsgLinkInvalid, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerRegisterLink.Register_link_is_invalid",
			Other: "the link is invalid or expired",
		},
	})

	linkDec, err := tools.DecryptBase64(link, app.conf.EncryptKey)
	if err != nil {
		page := &Page{
			ErrorMessage: &errMsgLinkInvalid,
		}
		logrus.Error(err)
		app.ExecuteTemplate(w, r, view, templateFile, page)
		return
	}

	linkUser := &models.RegisterUserLink{}

	err = json.Unmarshal([]byte(linkDec), linkUser)
	if err != nil {
		page := &Page{
			ErrorMessage: &errMsgLinkInvalid,
		}
		logrus.Error(err)
		app.ExecuteTemplate(w, r, view, templateFile, page)
		return
	}

	dbRegisterUser, err := repository.ReadRegisterUser(app.db, linkUser.ID)
	if err != nil {
		page := &Page{
			ErrorMessage: &errMsgLinkInvalid,
		}
		logrus.Error(err)
		app.ExecuteTemplate(w, r, view, templateFile, page)
		return
	}
	if dbRegisterUser.UpdatedAt != linkUser.CreatedAt {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Api.HandlerRegisterLink.Register_link_is_expired",
				Other: "the link is expired",
			},
		})
		page := &Page{
			ErrorMessage: &msg,
		}
		logrus.Error(errors.New("dbRegisterUser.UpdatedAt != linkUser.CreatedA"))
		app.ExecuteTemplate(w, r, view, templateFile, page)
		return
	}

	dbUser, err := repository.GetUserByEmail(app.db, dbRegisterUser.Email)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			msg, _ := localizer.Localize(&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID: DEFAULT_MESSAGE_COMMON_SERVER_ERROR,
				},
			})
			page := &Page{
				ErrorMessage: &msg,
			}
			logrus.Error(err)
			app.ExecuteTemplate(w, r, view, templateFile, page)
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
		page := &Page{
			ErrorMessage: &msg,
		}
		logrus.Error(errors.New(" dbUser != nil && dbUser.Email == dbRegisterUser.Email"))

		app.ExecuteTemplate(w, r, view, templateFile, page)
		return
	}

	dbUser = dbRegisterUser.ToUserModel()
	_, err = repository.CreateUser(app.db, dbUser)
	if err != nil {
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID: DEFAULT_MESSAGE_COMMON_SERVER_ERROR,
			},
		})
		page := &Page{
			ErrorMessage: &msg,
		}
		logrus.Error(err)

		app.ExecuteTemplate(w, r, view, templateFile, page)
		return
	}
	_ = repository.DeleteRegisterUser(app.db, dbRegisterUser.ID)

	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Api.HandlerRegisterLink.Message_user_created",
			Other: "successful, you can login now!",
		},
	})
	page := &Page{
		Message: &msg,
	}

	app.ExecuteTemplate(w, r, view, templateFile, page)

}

/*
func (app *App) RefreshLoginToken(w http.ResponseWriter, r *http.Request) {

	token := models.Token{}

	logrus.Debug("body: %d", r.Context())
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
