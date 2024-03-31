package app

import "net/http"

func (app *App) PageHome(res http.ResponseWriter, req *http.Request) {
	//app.printError(res, http.StatusInternalServerError, 200, nil, "")

	res.Write([]byte("Hallo1"))
}
