package echozap

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Log *zap.Logger
	// do not log 200-299 requests
	Skip2XX bool `json:"SKIP2XX"`
}

// ZapLogger is a middleware and zap to provide an "access log" like logging for each request.
func ZapLogger(config *Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)
			if err != nil {
				config.Log = config.Log.With(zap.Error(err))
				c.Error(err)
			}

			req := c.Request()
			res := c.Response()

			fields := []zapcore.Field{
				zap.String("remote_ip", c.RealIP()),
				zap.String("time", time.Since(start).String()),
				zap.String("host", req.Host),
				zap.String("request", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
				zap.Int("status", res.Status),
				zap.Int64("size", res.Size),
				zap.String("user_agent", req.UserAgent()),
			}

			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
				fields = append(fields, zap.String("request_id", id))
			}

			n := res.Status
			switch {
			case n >= 500:
				config.Log.Error("Server error", fields...)
			case n >= 400:
				config.Log.Warn("Client error", fields...)
			case n >= 300:
				config.Log.Info("Redirection", fields...)
			case n >= 200:
				if !config.Skip2XX {
					config.Log.Info("Success", fields...)
				}
			default:
				config.Log.Info("Success", fields...)
			}

			return nil
		}
	}
}
