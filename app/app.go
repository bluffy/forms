package app

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/bluffy/forms/config"
	"github.com/bluffy/forms/lang"
	"github.com/bluffy/forms/models"
	val "github.com/bluffy/forms/util/validator"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	appErr                    = "app error"
	appErrCreationFailure     = "error createn failure"
	appErrDataAccessFailure   = "data access failure"
	appErrJsonCreationFailure = "json creation failure"

	appErrDataCreationFailure = "data creation failure"
	appErrFormDecodingFailure = "form decoding failure"

	appErrDataUpdateFailure      = "data update failure"
	appErrFormErrResponseFailure = "form error response failure"
)

//go:embed version/*
var versionFS embed.FS

var version = "0.0.0"

type App struct {
	conf      *config.Config
	validator *validator.Validate
	lang      *lang.Lang
	db        *gorm.DB
	//userRestConf *clientcredentials.Config
	//	openIds map[string]*oauth2.Config
	//userClient *http.Client
	//amsClient  *http.Client
	user *models.UserDto
}

func New(
	conf *config.Config,
	validator *validator.Validate,
	lang *lang.Lang,
	db *gorm.DB,
) *App {

	return &App{
		conf:      conf,
		validator: validator,
		lang:      lang,
		//openIds: openIds,
		db: db,
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

func (app *App) SetUser(user models.UserDto) {
	app.user = &user
}

func (app *App) User() *models.UserDto {
	return app.user
}

func (app *App) Conf() *config.Config {
	return app.conf
}

func printError(app *App, w http.ResponseWriter, status int, msg string, err error) {
	if err != nil && app.conf.Debug {
		_, fn, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"func": fn,
			"line": fmt.Sprintf("%d", line),
		}).Error(err)
	} else {
		if err != nil {
			log.Warn(err)
		}
	}
	if msg != "" {
		log.Error(msg)
	}

	w.WriteHeader(status)
	message := ""

	if msg == "" {
		message = appErr
	} else {
		message = msg
	}
	errorObj := val.ErrorMsg(message)
	errorJson, err := json.Marshal(errorObj)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrCreationFailure)
		return
	}
	w.Write(errorJson)

}
func (a *App) printError(w http.ResponseWriter, status int, code int, err error, lang string) {
	if err != nil && a.conf.Debug {
		_, fn, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"func":    fn,
			"line":    fmt.Sprintf("%d", line),
			"code":    fmt.Sprintf("%d", code),
			"message": a.GetErrorByCode(code, a.lang.Log),
		}).Error(err)
	} else {
		if err != nil {
			log.Warn(err)
		}
	}

	var appError models.AppError
	var msg = a.GetErrorByCode(code, lang)
	appError.Error.Message = &msg
	if code != 0 {
		appError.Error.Code = &code
	}

	w.WriteHeader(status)

	errorJson, err := json.Marshal(appError)
	if err != nil {
		errMsg := a.GetErrorByCode(102, lang)
		log.Error("Error Message:" + a.GetErrorByCode(102, a.lang.Log))
		log.Error("Error Code" + strconv.Itoa(code))
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		if code != 0 {
			fmt.Fprintf(w, `{"error": {"message": "%s", "code": "%d"}}`, errMsg, code)
		} else {
			fmt.Fprintf(w, `{"error": {"message": "%s"}}`, errMsg)
		}
		return
	}
	w.Write(errorJson)

}

func PrintError(a *App, w http.ResponseWriter, status int, code int, err error, lang string) {
	a.printError(w, status, code, err, lang)

}

func (a *App) GetErrorByCode(code int, lang string) string {
	var message = "Unbekannter Fehler"

	var languageIdx = a.lang.Default
	if lang != "" {
		_, ok := a.lang.Region[lang]
		if !ok {
			languageIdx = lang
		}
	}

	region, ok := a.lang.Region[languageIdx]

	var errMsg = ""
	if ok {
		errMsg, ok = region.Error[code]
		if ok {
			message = errMsg
		} else {
			region, _ = a.lang.Region[a.lang.Default]
			errMsg, ok = region.Error[code]
			if ok {
				message = errMsg
			}
		}
	}

	return message

}

/*
func (a *App) printErrorByCode(w http.ResponseWriter, status int, code int, err error, lang string) {
	a.printError(w, status, code, a.GetErrorByCode(code, lang), err, lang)
}
*/

func (app *App) checkForm(form interface{}, w http.ResponseWriter, r *http.Request) (stop bool) {
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return true
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := val.ToErrResponse(err, nil)
		if resp == nil {
			printError(app, w, http.StatusInternalServerError, appErrFormErrResponseFailure, err)
			return true
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			printError(app, w, http.StatusInternalServerError, appErrJsonCreationFailure, err)
			return true
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return true
	}

	return false
}
