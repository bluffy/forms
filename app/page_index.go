package app

import (
	"net/http"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
)

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
// @Router       /bl-api/page/v1/ [get]
func (app *App) PageIndex(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")
	logrus.Println("#################")

	localizer := req.Context().Value(LocalizerKey{}).(*i18n.Localizer)
	msg, _ := localizer.Localize(&i18n.LocalizeConfig{

		PluralCount: 2,

		DefaultMessage: &i18n.Message{
			ID:          "HelloWorld6",
			Many:        "hallo viele ",
			Few:         "ein ppar  !",
			Other:       "Hell  ",
			Description: "beschreibung",
			One:         "Hallo 1 {{.PluralCount}} ",
		},
	})

	logrus.Println("###############")
	logrus.Println(msg)

	//mysess := session.GetSession(req)

	//logrus.Printf("%+v\n", mysess.Get("name"))
	//	func (m *Manager) Start(resp, req *http.Request) (RawStore, error)

	//logrus.Printf("%+v\n", mysess.)
	res.Write([]byte("{\"OK\": true}"))

}
