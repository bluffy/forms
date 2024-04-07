package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"goapp/config"
	"goapp/lang"
	"goapp/models"
	"goapp/repository"

	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:embed version/*
var versionFS embed.FS
var version = "0.0.0"

type ErrResponse struct {
	//Errors []string `json:"errors"`
	Error struct {
		Message *string                 `json:"message,omitempty"`
		Fields  *map[string]interface{} `json:"fields,omitempty"`
	} `json:"error"`
}

type PageResponse struct {
	Message *string `json:"message,omitempty"`
}

type App struct {
	validator *validator.Validate
	lang      *lang.Lang
	db        *gorm.DB
	conf      *config.Config
	bundle    *i18n.Bundle
}

func New(
	validator *validator.Validate,
	lang *lang.Lang,
	db *gorm.DB,
	config *config.Config,
	bundle *i18n.Bundle,
) *App {

	return &App{
		validator: validator,
		lang:      lang,
		db:        db,
		conf:      config,
		bundle:    bundle,
	}
}

func init() {
	data, _ := versionFS.ReadFile("version/VERSION")
	if data != nil {
		version = string(data)
	}
}

func GetVersion() string {
	return version
}
func (a *App) Conf() config.Config {
	return *a.conf
}

func (a *App) GetLocalizer(r *http.Request, lang string) *i18n.Localizer {
	accept := r.Header.Get("Accept-Language")
	return i18n.NewLocalizer(a.bundle, lang, accept)

}
func (a *App) JsonError(w http.ResponseWriter, status int, publicMessage string, err error, doLog bool, optionalLoggingMessage string) {
	commonError := "Error on Server"
	_, fn, line, _ := runtime.Caller(1)
	if doLog || a.conf.Debug {

		log.WithFields(log.Fields{
			"jsonError":      true,
			"func":           fn,
			"line":           fmt.Sprintf("%d", line),
			"publicMessage":  publicMessage,
			"loggingMessage": optionalLoggingMessage,
			"httpsStatus":    status,
		}).Error(err)

	}

	var appError ErrResponse
	if publicMessage == "" {
		msg := a.GetLocale("").Text.Error__commen_server_error
		appError.Error.Message = &msg
	} else {
		appError.Error.Message = &publicMessage
	}

	w.WriteHeader(status)

	errorJson, err := json.Marshal(appError)
	if err != nil {
		log.WithFields(log.Fields{
			"func":           fn,
			"line":           fmt.Sprintf("%d", line),
			"publicMessage":  publicMessage,
			"loggingMessage": optionalLoggingMessage,
		}).Error(err)
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, commonError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(errorJson)
}

func (a *App) errMessage(publicMessage string, err error, doLog bool, optionalLoggingMessage string) string {

	if doLog || a.conf.Debug {
		_, fn, line, _ := runtime.Caller(1)

		//log.Info("####### Error (TODOD write to sql)")
		log.WithFields(log.Fields{
			"errMessage":     true,
			"func":           fn,
			"line":           fmt.Sprintf("%d", line),
			"publicMessage":  publicMessage,
			"loggingMessage": optionalLoggingMessage,
		}).Error(err)

	}

	if publicMessage == "" {
		msg := a.GetLocale("").Text.Error__commen_server_error
		return msg
	}

	return publicMessage

}

func (a *App) GetLocale(lang string) *lang.Locale {

	locale, ok := a.lang.Locale[lang]
	if !ok {
		locale = *a.lang.DefaultLocale
	}

	return &locale
}

func (app *App) formErrors(lang string, err error, msg *string) *ErrResponse {

	if fieldErrors, ok := err.(validator.ValidationErrors); ok {

		resp := ErrResponse{}
		resp.Error.Message = msg
		fields := make(map[string]interface{})

		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Required)
			case "max":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Max, err.Param())
			case "min":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Min, err.Param())
			case "url":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Url)
			case "email":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Email)
			case "alpha_space":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Alpha_space)
			case "date":
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Date)
			default:
				fields[err.Field()] = fmt.Sprintf(app.GetLocale(lang).Validator.Default, err.Tag())
			}
		}
		if len(fields) > 0 {
			resp.Error.Fields = &fields
		}
		return &resp
	}
	return nil
}

func (app *App) checkForm(lang string, form interface{}, w http.ResponseWriter, r *http.Request, errorMessage *string) (stop bool) {
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__form_response_error, err, false, "")
		return true
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := app.formErrors(lang, err, errorMessage)
		if resp == nil {
			app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__form_response_error, err, false, "")
			return true
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			app.JsonError(w, http.StatusInternalServerError, app.GetLocale("").Text.Error__json_create, err, false, "")
			return true
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return true
	}

	return false
}

func (a *App) sendMail(adhoc bool, mailSubject string, mailText string, mailHtml string, to string) error {
	var mail models.Mail
	mail.Status = 0
	mail.Text = &mailText
	if mailHtml != "" {
		mail.Html = &mailHtml
	}
	mail.Recipient = to
	mail.Subject = mailSubject
	mail.Sender = a.conf.Smtp.Sender
	if adhoc {
		mail.Status = models.SEND_STATUS_WAITING
	}
	dbMail, err := repository.CreateMail(a.db, &mail)
	if err != nil {
		//a.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__database_error, "Error__database_error in CreateMail", err)
		return err
	}

	if adhoc {
		serviceMail := dbMail.ToServiceMail()
		logMsg, err := serviceMail.SendMail(&a.conf.Smtp)

		if err != nil {
			mailError := fmt.Sprintf("%v", err)
			dbMail.ErrorMessage = logMsg
			dbMail.Error = &mailError
			dbMail.Status = models.SEND_STATUS_ERROR
		} else {
			dbMail.Status = models.SEND_STATUS_SENT
			now := time.Now()
			dbMail.SendAt = &now
		}
		_, errDB := repository.UpdateMail(a.db, dbMail)
		if errDB != nil {
			//app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__database_error, "Error__database_error in UpdateMail", err)
			return err
		}
	}

	return err
}
