package app

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

func (app *App) HanlderHealth(res http.ResponseWriter, req *http.Request) {

	log.Printf("--- m dump:\n%s\n\n", req.Header.Get("Accept-Language"))

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	val, ok := app.lang.Locale["en"]
	log.Println(ok)
	if ok {
		log.Println(val.Text.Welcome)
	}
	d, err := yaml.Marshal(&val)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("--- m dump:\n%s\n\n", string(d))
	/*
		log.Println("TEST")

		val, ok := app.lang.Locale["en"]
		log.Println(ok)
		if ok {
			log.Println(val.Text.Welcome)
		}

		d, err := yaml.Marshal(&val)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Printf("--- m dump:\n%s\n\n", string(d))
	*/

	/*
		var err error
		if !config.Conf.UseDad {
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
