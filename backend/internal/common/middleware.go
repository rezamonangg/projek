package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

type httpLog struct {
	Method     string `json:"method"`
	URL        string `json:"url"`
	Status     int    `json:"status"`
	LatencyMs  int64  `json:"latency_ms"`
	RequestID  string `json:"request_id,omitempty"`
	RemoteAddr string `json:"remote_addr,omitempty"`
}

type requestLogEntry struct {
	Timestamp string     `json:"timestamp"`
	Level     string     `json:"level"`
	AppName   string     `json:"appname"`
	Message   string     `json:"message"`
	Context   LogContext `json:"context"`
	HTTP      httpLog    `json:"http"`
}

func RequestLogger(_ zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()
			defer func() {
				entry := requestLogEntry{
					Timestamp: time.Now().Format(time.RFC3339),
					Level:     "info",
					AppName:   "projek",
					Message:   "request completed",
					Context:   GetLogContext(r.Context()),
					HTTP: httpLog{
						Method:     r.Method,
						URL:        r.URL.String(),
						Status:     ww.Status(),
						LatencyMs:  time.Since(start).Microseconds(),
						RequestID:  middleware.GetReqID(r.Context()),
						RemoteAddr: r.RemoteAddr,
					},
				}

				data, _ := json.Marshal(entry)
				fmt.Println(string(data))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
