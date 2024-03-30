package app

import (
	"net/http"
)

func (app *App) HanlderHealth(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	/*
		var err error
		if !app.conf.UseDad {
			err = oracle.Ping()
		}

		if err != nil {
			res.WriteHeader(http.StatusBadGateway)
			res.Write([]byte("NOT OK, DATABASE NOT REACHABLE"))
			return
		}
	*/

	// Write the status code using w.WriteHeader
	res.WriteHeader(http.StatusOK)

	// Write the body text using w.Write
	res.Write([]byte("OK"))
}
