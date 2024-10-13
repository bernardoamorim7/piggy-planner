package middlewares

import (
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type RequestLog struct {
	Method     string
	URL        string
	RemoteAddr string
	Timestamp  time.Time
}

var (
	requestLogs []RequestLog
	mu          sync.Mutex
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log := RequestLog{
				Method:     c.Request().Method,
				URL:        c.Request().URL.String(),
				RemoteAddr: c.Request().RemoteAddr,
				Timestamp:  time.Now(),
			}

			mu.Lock()
			if log.URL != "/api/requests" && log.URL != "/api/requests/history" {
				requestLogs = append([]RequestLog{log}, requestLogs...) // Insert at the beginning
			}
			mu.Unlock()

			return next(c)
		}
	}
}

func GetRequestLogs() []RequestLog {
	mu.Lock()
	defer mu.Unlock()
	return requestLogs
}
