package aspects

import (
	"bufio"
	"bytes"
	"fmt"
	"go.uber.org/zap/zapcore"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type LoggingMiddlewareConfig struct {
	IncludeRequestBodies  bool
	IncludeResponseBodies bool
}

func NewLoggingMiddleware(logger *zap.Logger, config *LoggingMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.String(), "healthz") {
				return next(c)
			}

			if strings.Contains(c.Request().URL.String(), "ready") {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Request
			reqBody := []byte{}
			if config.IncludeRequestBodies {
				if req.Body != nil { // Read
					reqBody, _ = ioutil.ReadAll(req.Body)
				}
				req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Cleanup
			}

			// Response
			respBody := new(bytes.Buffer)
			if config.IncludeResponseBodies {
				mw := io.MultiWriter(res.Writer, respBody)
				writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
				res.Writer = writer
			}

			if err := next(c); err != nil {
				c.Error(err)
			}

			logEntry := logger

			if config.IncludeRequestBodies {
				logEntry = logEntry.With(zap.Field{
					Key:    "requestBody",
					Type:   zapcore.StringType,
					String: string(reqBody),
				})
			}

			if config.IncludeResponseBodies {
				logEntry = logEntry.With(zap.Field{
					Key:    "responseBody",
					Type:   zapcore.StringType,
					String: respBody.String(),
				})
			}

			logEntry.Info(fmt.Sprintf("Request %s %s finished in %dms", req.Method, req.URL, time.Since(start).Milliseconds()))

			return nil
		}
	}
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *bodyDumpResponseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
