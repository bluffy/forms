package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"

	"goapp/config"
	"goapp/lang"

	"github.com/go-playground/validator/v10"
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

type App struct {
	validator *validator.Validate
	lang      *lang.Lang
	db        *gorm.DB
	conf      *config.Config
}

func New(
	validator *validator.Validate,
	lang *lang.Lang,
	db *gorm.DB,
	config *config.Config,
) *App {

	return &App{
		validator: validator,
		lang:      lang,
		db:        db,
		conf:      config,
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

func (a *App) JsonError(w http.ResponseWriter, status int, publicMessage string, loggingMessage string, err error) {
	commonError := "Error on Server"
	_, fn, line, _ := runtime.Caller(1)
	if err != nil && a.conf.Debug {
		log.WithFields(log.Fields{
			"jsonError":      true,
			"func":           fn,
			"line":           fmt.Sprintf("%d", line),
			"publicMessage":  publicMessage,
			"loggingMessage": loggingMessage,
			"httpsStatus":    status,
		}).Error(err)
	} else {
		if err != nil {
			log.Warn(err)
		}
	}

	if loggingMessage != "" {
		log.Info("####### Error (TODOD write to sql)")
		log.WithFields(log.Fields{
			"jsonError":      true,
			"func":           fn,
			"line":           fmt.Sprintf("%d", line),
			"publicMessage":  publicMessage,
			"loggingMessage": loggingMessage,
			"httpsStatus":    status,
		}).Info(err)
		log.Info("#######")
	}

	var appError ErrResponse
	if publicMessage == "" {
		msg := "Error on Server"
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
			"loggingMessage": loggingMessage,
		}).Error(err)
		fmt.Fprintf(w, `{"error": {"message": "%s"}}`, commonError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(errorJson)
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
		app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__form_response_error, "Error__form_response_error", err)
		return true
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := app.formErrors(lang, err, errorMessage)
		if resp == nil {
			app.JsonError(w, http.StatusUnprocessableEntity, app.GetLocale("").Text.Error__form_response_error, "Error__form_response_error", err)
			return true
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			app.JsonError(w, http.StatusInternalServerError, app.GetLocale("").Text.Error__json_create, "Error__json_create", err)
			return true
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return true
	}

	return false
}
