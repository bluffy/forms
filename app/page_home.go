package app

import (
	"goapp/config"
	"goapp/util/tools"
	"log"
	"net/http"
)

func (app *App) PageHome(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	session, err := req.Cookie("session")
	if err != nil {

		log.Print(err)
		res.Write([]byte("fehler"))
		return

	}
	id, err := tools.DecryptBase64(session.Value, config.Conf.EncryptKey)
	if err != nil {

		log.Print(err)
		res.Write([]byte("fehler"))
		return

	}
	log.Print(session)

	res.Write([]byte("Hallo:" + id + ": session"))
}
