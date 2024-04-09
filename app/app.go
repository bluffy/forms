package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"runtime"
	"time"

	"goapp/config"
	"goapp/models"
	"goapp/repository"
	"goapp/util/tools"

	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:embed version/*
var versionFS embed.FS
var version = "0.0.0"

type ContextUserIDKey struct{}
type ContextLocalizerKey struct{}
type ContextSessionStoreKey struct{}

type ErrResponse struct {
	//Errors []string `json:"errors"`
	Error struct {
		Message *string                 `json:"message,omitempty"`
		Fields  *map[string]interface{} `json:"fields,omitempty"`
	} `json:"error"`
}

type ApiPageResponse struct {
	Message *string `json:"message,omitempty"`
}

type Page struct {
	ErrorMessage *string
	Message      *string
}
type PageError struct {
	Message string
}

type App struct {
	validator  *validator.Validate
	db         *gorm.DB
	conf       *config.Config
	bundle     *i18n.Bundle
	templateFS *embed.FS
}

func New(
	validator *validator.Validate,
	db *gorm.DB,
	config *config.Config,
	bundle *i18n.Bundle,
	templateFS *embed.FS,
) *App {
	loadDefaultMessages(bundle)
	return &App{
		validator:  validator,
		db:         db,
		conf:       config,
		bundle:     bundle,
		templateFS: templateFS,
	}
}

const DEFAULT_ERROR = "ServerError"

const LOCALE_MSG_ID__COMMON_SERVER_ERROR = "Server.CommonError"

func loadDefaultMessages(bundle *i18n.Bundle) {
	localizer := i18n.NewLocalizer(bundle)

	_, _ = localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    LOCALE_MSG_ID__COMMON_SERVER_ERROR,
			Other: "Error on Server",
		},
	})

}

func GetDefaultMessageError(r *http.Request) string {
	localizer := GetLocalizer(r)
	if localizer == nil {
		return DEFAULT_ERROR
	}
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: LOCALE_MSG_ID__COMMON_SERVER_ERROR,
		},
	})
	return msg
}
func GetDefaultMessage(localizer *i18n.Localizer, id string) string {
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: id,
		},
	})
	return msg
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

func (a *App) GetBundle() *i18n.Bundle {
	return a.bundle
}

func GetLocalizer(r *http.Request) *i18n.Localizer {
	return r.Context().Value(ContextLocalizerKey{}).(*i18n.Localizer)
}

func (a *App) JsonErrorMessage(locConfig *i18n.LocalizeConfig, r *http.Request, w http.ResponseWriter, status int, err error, doLog bool, optionalLoggingMessage string) {
	publicMessage, _ := GetLocalizer(r).Localize(locConfig)

	commonError := "Error on Server"
	_, fn, line, _ := runtime.Caller(1)
	if doLog || a.conf.Debug {

		logrus.WithFields(logrus.Fields{
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
		msg := GetDefaultMessage(GetLocalizer(r), LOCALE_MSG_ID__COMMON_SERVER_ERROR)
		appError.Error.Message = &msg
	} else {
		appError.Error.Message = &publicMessage
	}

	w.WriteHeader(status)

	errorJson, err := json.Marshal(appError)
	if err != nil {
		logrus.WithFields(logrus.Fields{
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

func (a *App) JsonError(r *http.Request, w http.ResponseWriter, status int, publicMessage string) {

	var appError ErrResponse
	if publicMessage == "" {
		commonError := GetDefaultMessageError(r)
		appError.Error.Message = &commonError
	} else {
		appError.Error.Message = &publicMessage
	}
	w.WriteHeader(status)

	errorJson, err := json.Marshal(appError)
	if err != nil {
		commonError := GetDefaultMessageError(r)
		a.ServerLogByRequest(r, err, commonError, true, "json.Marshal(appError)")
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, commonError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(errorJson)
}

func (a *App) Render(w io.Writer, p Page) error {
	templ, err := template.ParseFS(a.templateFS, "templates/*")
	if err != nil {
		return err
	}

	if err := templ.Execute(w, p); err != nil {
		return err
	}

	return nil
}

func (a *App) serverLog(err error, logMessage string) {
	_, fn, line, _ := runtime.Caller(1)

	fields := logrus.Fields{
		"func":       fn,
		"line":       fmt.Sprintf("%d", line),
		"logMessage": logMessage,
	}

	logrus.WithFields(fields).Error(err)
}

func (a *App) ServerLogByRequest(r *http.Request, err error, publicMessage string, doLog bool, logMessage string) {
	if doLog || a.conf.Debug {
		_, fn, line, _ := runtime.Caller(1)

		fields := logrus.Fields{
			"func":       fn,
			"line":       fmt.Sprintf("%d", line),
			"logMessage": logMessage,
			"PublicIP":   tools.ReadUserIP(r),
			"Host":       r.Host,
			"ApiPath":    r.URL.Path,
		}
		if publicMessage != "" {
			fields["publicMessage"] = publicMessage
		}
		if publicMessage != "" {
			fields["logMessage"] = logMessage
		}
		if a.conf.Debug {
			fields["doLog"] = doLog
		}

		logrus.WithFields(fields).Error(err)
	}
}

func (a *App) RenderError(r *http.Request, w io.Writer, p *PageError) {

	if p == nil {
		p = &PageError{
			Message: DEFAULT_ERROR,
		}
	}

	templ, tmpErr := template.ParseFS(a.templateFS, "templates/page/*")

	if tmpErr != nil {
		a.ServerLogByRequest(r, tmpErr, "", true, "")
		w.Write([]byte(GetDefaultMessageError(r)))
	}

	if tmpErr := templ.ExecuteTemplate(w, "error.gohtml", p); tmpErr != nil {
		a.ServerLogByRequest(r, tmpErr, "", true, "")
		w.Write([]byte(GetDefaultMessageError(r)))
	}

}
func (a *App) RenderPage(r *http.Request, w io.Writer, p *Page) {

	if p == nil {
		msg := GetDefaultMessageError(r)
		a.ServerLogByRequest(r, nil, msg, true, "p *Page is null")
		w.Write([]byte(msg))
	}

	templ, tmpErr := template.ParseFS(a.templateFS, "templates/page/*")

	if tmpErr != nil {
		msg := GetDefaultMessageError(r)
		a.ServerLogByRequest(r, tmpErr, msg, true, "")
		w.Write([]byte(msg))
	}

	if tmpErr := templ.ExecuteTemplate(w, "error.gohtml", p); tmpErr != nil {
		msg := GetDefaultMessageError(r)
		a.ServerLogByRequest(r, tmpErr, msg, true, "")
		w.Write([]byte(msg))
	}

}

func (a *App) ExecuteTemplate(w http.ResponseWriter, r *http.Request, view *template.Template, templateFile string, data any) {
	err := view.ExecuteTemplate(w, templateFile, data)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		msg := GetDefaultMessage(GetLocalizer(r), LOCALE_MSG_ID__COMMON_SERVER_ERROR)
		logrus.WithFields(logrus.Fields{
			"func":           fn,
			"line":           fmt.Sprintf("%d", line),
			"publicMessage":  msg,
			"loggingMessage": "view.ExecuteTemplate",
		}).Error(err)
		w.Write([]byte(msg))
	}

}

func (app *App) formErrors(localizer *i18n.Localizer, err error, msg *string) *ErrResponse {

	if fieldErrors, ok := err.(validator.ValidationErrors); ok {

		resp := ErrResponse{}
		resp.Error.Message = msg
		fields := make(map[string]interface{})

		for _, err := range fieldErrors {
			msg := ""
			switch err.Tag() {
			case "required":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.Required",
						Other: "is required",
					},
				})
			case "max":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					TemplateData: map[string]string{
						"Chars": err.Param(),
					},
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.Max",
						Other: "must be a maximum of {{.Chars}} in length",
					},
				})
			case "min":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					TemplateData: map[string]string{
						"Chars": err.Param(),
					},
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.Min",
						Other: "must be a minimum of {{.Chars}} in length",
					},
				})
			case "url":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.URL",
						Other: "must be a valid URL",
					},
				})
			case "email":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.Email",
						Other: "must be a valid email addres",
					},
				})
			case "alpha_space":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.AlphaSpace",
						Other: "can only contain alphabetic and space characters",
					},
				})
			case "date":
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.Date",
						Other: "imust be a valid date",
					},
				})
			default:
				msg, _ = localizer.Localize(&i18n.LocalizeConfig{
					TemplateData: map[string]string{
						"Tag": err.Tag(),
					},
					DefaultMessage: &i18n.Message{
						ID:    "FormsValidator.Default",
						Other: "something wrong: {{.Tag}}",
					},
				})
			}
			fields[err.Field()] = msg
		}
		if len(fields) > 0 {
			resp.Error.Fields = &fields
		}
		return &resp
	}
	return nil
}

func (app *App) checkForm(localizer *i18n.Localizer, form interface{}, w http.ResponseWriter, r *http.Request, errorMessage *string) (stop bool) {
	msgFromError, _ := localizer.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "FormsValidator.FormResponseError",
			Other: "form has an error",
		},
	})

	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.ServerLogByRequest(r, err, msgFromError, false, "json.NewDecoder(r.Body).Decode(form)")
		app.JsonError(r, w, http.StatusUnprocessableEntity, msgFromError)
		return true
	}

	if err := app.validator.Struct(form); err != nil {
		logrus.Warn(err)
		resp := app.formErrors(localizer, err, errorMessage)
		if resp == nil {
			app.ServerLogByRequest(r, err, msgFromError, false, "app.formErrors(localizer, err, errorMessage)")
			app.JsonError(r, w, http.StatusUnprocessableEntity, msgFromError)
			return true
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			app.ServerLogByRequest(r, err, msgFromError, true, "json.Marshal(resp)")
			app.JsonError(r, w, http.StatusInternalServerError, "")
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
			return err
		}
	}

	return err
}
