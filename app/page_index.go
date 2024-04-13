package app

import (
	"net/http"

	"gitea.com/go-chi/session"
	"github.com/sirupsen/logrus"
)

func (app *App) PageIndex(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	//localizer := req.Context().Value(ContextLocalizerKey{}).(*i18n.Localizer)
	/*
		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			PluralCount: 2,
			DefaultMessage: &i18n.Message{
				ID:          "HelloWorld7",
				Many:        "hallo viele ",
				Few:         "ein ppar  !",
				Other:       "Hell 33 ",
				Description: "beschreibung",
				One:         "Hallo 1 {{.PluralCount}} ",
			},
		})
	*/
	msg := "TEST"

	//msg := "TEST"
	sess := session.GetSession(req)
	logrus.Info(sess.Get("test-session"))

	logrus.Info(sess.Get("user_id"))

	logrus.Println("###############")
	logrus.Println(msg)

	//mysess := session.GetSession(req)

	//logrus.Printf("%+v\n", mysess.Get("name"))
	//	func (m *Manager) Start(resp, req *http.Request) (RawStore, error)

	//logrus.Printf("%+v\n", mysess.)
	res.Write([]byte("{\"OK\": \"" + msg + "\"}"))

}
