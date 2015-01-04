package wrappers

import (
	"net/http"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

func logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLogger := log.New(
			"name", name,
			"method", r.Method,
			"uri", r.RequestURI,
		)
		start := time.Now()

		inner.ServeHTTP(w, r)

		requestLogger.Info("request complete", "duration", time.Since(start))
	})
}

func cors(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods",
			"HEAD,OPTIONS,GET,PUT,PATCH,POST,DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept,Content-Type,Content-Length,Accept-Encoding,Authorization,"+
				"X-Requested-With",
		)
		if r.Method == "OPTIONS" {
			return
		}
		inner.ServeHTTP(w, r)
	})
}

// Defaults adds CORS headers and request logging to every request.
func Defaults(inner http.HandlerFunc, name string) http.Handler {
	return logger(cors(http.HandlerFunc(inner)), name)
}
