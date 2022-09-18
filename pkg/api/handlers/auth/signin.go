package auth

import (
	"net/http"

	"apart-deal-api/pkg/api/auth"

	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo/v4"

	apiErr "apart-deal-api/pkg/api/aspects/errors"
	validation "github.com/go-ozzo/ozzo-validation"
	oas "gitlab.com/apart-deals/openapi/go/api"
)

func validateSignIn(payload *oas.SignIn) error {
	return validation.ValidateStruct(
		payload,
		validation.Field(&payload.Email, validation.Required, is.Email, validation.Length(3, 50)),
		validation.Field(&payload.Password, validation.Required, validation.Length(4, 10)),
	)
}

type SignInHandler struct {
	authSvc *auth.AuthenticationService
}

func NewSignInHandler(authSvc *auth.AuthenticationService) *SignInHandler {
	return &SignInHandler{
		authSvc: authSvc,
	}
}

func (h *SignInHandler) Handle(eCtx echo.Context) error {
	payload := &oas.SignIn{}

	if err := eCtx.Bind(payload); err != nil {
		return err
	}

	if err := validateSignIn(payload); err != nil {
		return apiErr.NewMultipleValidationInputError(err)
	}

	tokenString, err := h.authSvc.Auth(eCtx.Request().Context(), payload)
	if err != nil {
		return mapError(err)
	}

	return eCtx.JSON(http.StatusOK, oas.SignedIn{
		Token: tokenString,
	})
}
