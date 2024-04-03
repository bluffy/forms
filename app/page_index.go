package app

import (
	"net/http"
)

func (app *App) PageIndex(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	//mysess := session.GetSession(req)

	//log.Printf("%+v\n", mysess.Get("name"))
	//	func (m *Manager) Start(resp, req *http.Request) (RawStore, error)

	//log.Printf("%+v\n", mysess.)
	res.Write([]byte("{\"OK\": true}"))

}
