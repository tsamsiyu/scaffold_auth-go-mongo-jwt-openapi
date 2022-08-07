package aspects

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func NewErrorHandler(logger *zap.Logger) echo.HTTPErrorHandler {
	return func(err error, context echo.Context) {
		httpErr, ok := err.(*echo.HTTPError)
		if ok {
			msgStr, isMsgStr := httpErr.Message.(string)
			if isMsgStr {
				_ = context.JSON(httpErr.Code, map[string]string{
					"message": msgStr,
				})
			} else {
				_ = context.JSON(httpErr.Code, httpErr.Message)
			}

			return
		}

		logger.Error(
			fmt.Sprintf("Unhandled error at [%s %s]: %s",
				context.Request().Method,
				context.Request().URL,
				errors.WithStack(err),
			),
		)

		_ = context.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "Something went wrong",
		})
	}
}
