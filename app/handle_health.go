package app

import (
	"goapp/repository"
	"net/http"
)

func (app *App) HanlderHealth(res http.ResponseWriter, req *http.Request) {

	if (repository.CountAllUsers(app.db)) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusOK)

	res.Write([]byte("OK"))
}
