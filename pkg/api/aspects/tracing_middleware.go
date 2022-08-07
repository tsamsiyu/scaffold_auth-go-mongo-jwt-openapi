package aspects

import (
	"strings"

	"apart-deal-api/pkg/tracing"

	"github.com/labstack/echo/v4"
)

func NewTracingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.Path, "healthz") {
				return next(c)
			}

			if strings.Contains(c.Request().URL.Path, "ready") {
				return next(c)
			}

			if traceID := c.Request().Header.Get("X-Trace-Id"); traceID != "" {
				newCtx := tracing.WithTraceID(c.Request().Context(), traceID)
				c.SetRequest(c.Request().WithContext(newCtx))
			}

			return next(c)
		}
	}
}
