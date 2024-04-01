package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"

	"goapp/util/tools"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

func Logger(category string, logger logrus.FieldLogger) func(h http.Handler) http.Handler {

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				fields := logrus.Fields{
					"status_code":      ww.Status(),
					"bytes":            ww.BytesWritten(),
					"duration":         int64(time.Since(t1)),
					"duration_display": time.Since(t1).String(),
					"category":         category,
					"remote_ip":        remoteIP,
					"proto":            r.Proto,
					"method":           r.Method,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}

				uri := r.RequestURI
				if strings.HasPrefix(uri, "/api/v1/f/") {
					uri = uri[0:tools.Min(100, len(uri))]
					if len(r.RequestURI) >= 100 {
						uri = uri + "..."
					}
				} else {
					uri = uri[0:tools.Min(255, len(uri))]
					if len(r.RequestURI) >= 255 {
						uri = uri + "..."
					}
				}

				if logger != nil {
					logger.WithFields(fields).Infof("%s://%s%s", scheme, r.Host, uri)
				} else {
					logrus.WithFields(fields).Infof("%s://%s%s", scheme, r.Host, uri)
				}
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
