package auth

import (
	"net/http"

	"apart-deal-api/pkg/api/auth"

	"github.com/labstack/echo/v4"

	apiErr "apart-deal-api/pkg/api/aspects/errors"
	validation "github.com/go-ozzo/ozzo-validation"
	oas "gitlab.com/apart-deals/openapi/go/api"
)

func validateRefreshAuthToken(payload *oas.RefreshAuthToken) error {
	return validation.ValidateStruct(
		payload,
		validation.Field(&payload.AuthToken, validation.Required),
		validation.Field(&payload.RefreshToken, validation.Required),
	)
}

type SignInRefreshHandler struct {
	authSvc *auth.AuthenticationService
}

func (h *SignInRefreshHandler) Refresh(eCtx echo.Context) error {
	payload := &oas.RefreshAuthToken{}

	if err := eCtx.Bind(payload); err != nil {
		return err
	}

	if err := validateRefreshAuthToken(payload); err != nil {
		return apiErr.NewMultipleValidationInputError(err)
	}

	t, err := h.authSvc.RefreshToken(eCtx.Request().Context(), payload)
	if err != nil {
		return mapError(err)
	}

	return eCtx.JSON(http.StatusOK, t)
}
