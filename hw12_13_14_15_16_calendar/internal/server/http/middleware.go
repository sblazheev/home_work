package internalhttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
)

func loggingMiddleware(next http.Handler, logger common.LoggerInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		rw := NewLoggingResponseWriter(w)
		next.ServeHTTP(rw, r)

		l := time.Since(s)
		logger.Info("REQUEST API",
			"data",
			struct {
				IP        string
				Date      time.Time
				Path      string
				Proto     string
				Method    string
				UserAgent string
				Status    int
				Latency   int
			}{
				IP:        strings.Split(r.RemoteAddr, ":")[0],
				Date:      s,
				Path:      r.URL.Path,
				Proto:     r.Proto,
				Method:    r.Method,
				UserAgent: r.UserAgent(),
				Status:    rw.statusCode,
				Latency:   int(l.Milliseconds()),
			})
	})
}
