package aspects

import (
	"fmt"
	"net/http"

	apiErr "apart-deal-api/pkg/api/aspects/errors"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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

		if _, ok := err.(*apiErr.ValidationError); ok {
			_ = context.JSON(http.StatusBadRequest, err)
			return
		}

		if _, ok := err.(*apiErr.ConflictError); ok {
			_ = context.JSON(http.StatusConflict, err)
			return
		}

		if _, ok := err.(*apiErr.NotFoundError); ok {
			_ = context.JSON(http.StatusNotFound, err)
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
