package app

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed client/*
var client embed.FS
var contentFS, _ = fs.Sub(client, "client")

type (
	// FallbackResponseWriter wraps an http.Requesthandler and surpresses
	// a 404 status code. In such case a given local file will be served.
	FallbackResponseWriter struct {
		WrappedResponseWriter http.ResponseWriter
		FileNotFound          bool
	}
)

func (frw *FallbackResponseWriter) WriteHeader(statusCode int) {
	//log.Printf("INFO: WriteHeader called with code %d\n", statusCode)
	if statusCode == http.StatusNotFound {
		//log.Printf("INFO: Setting FileNotFound flag\n")
		frw.FileNotFound = true
		return
	}
	frw.WrappedResponseWriter.WriteHeader(statusCode)
}

// Header returns the header of the wrapped response writer
func (frw *FallbackResponseWriter) Header() http.Header {
	return frw.WrappedResponseWriter.Header()
}

// Write sends bytes to wrapped response writer, in case of FileNotFound
// It surpresses further writes (concealing the fact though)
func (frw *FallbackResponseWriter) Write(b []byte) (int, error) {
	if frw.FileNotFound {
		return len(b), nil
	}
	return frw.WrappedResponseWriter.Write(b)
}

func (app *App) HandleClient(w http.ResponseWriter, r *http.Request) {

	frw := FallbackResponseWriter{
		WrappedResponseWriter: w,
		FileNotFound:          false,
	}

	http.FileServer(http.FS(contentFS)).ServeHTTP(&frw, r)
	if frw.FileNotFound {
		b, _ := client.ReadFile("client" + "/index.html")
		w.Header().Set("Content-Type", "text/html")
		w.Write(b)
		return
	}

}
