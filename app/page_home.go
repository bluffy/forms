package app

import (
	"goapp/config"
	"goapp/util/tools"
	"log"
	"net/http"

	"gitea.com/go-chi/session"
)

func (app *App) PageHome(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	mysess := session.GetSession(req)

	log.Printf("%v", mysess)

	sess, err := req.Cookie("session")
	if err != nil {

		log.Print(err)
		res.Write([]byte("fehler"))
		return

	}
	id, err := tools.DecryptBase64(sess.Value, config.Conf.EncryptKey)
	if err != nil {

		log.Print(err)
		res.Write([]byte("fehler"))
		return

	}
	log.Print(sess)

	res.Write([]byte("Hallo:" + id + ": sess"))
}
