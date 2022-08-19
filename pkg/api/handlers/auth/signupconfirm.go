package auth

import (
	"net/http"

	"apart-deal-api/pkg/api/aspects/errors"

	"github.com/labstack/echo/v4"

	authDomain "apart-deal-api/pkg/domain/auth"

	validation "github.com/go-ozzo/ozzo-validation"
	oas "gitlab.com/apart-deals/openapi/go/api"
)

func validateSignUpConfirm(input *oas.ConfirmSignUp) error {
	return validation.ValidateStruct(
		input,
		validation.Field(&input.Code, validation.Required),
		validation.Field(&input.Token, validation.Required),
	)
}

type SignUpConfirmHandler struct {
	confirmSignUpSvc *authDomain.ConfirmSignUpService
}

func NewSignUpConfirmHandler(confirmSignUpSvc *authDomain.ConfirmSignUpService) *SignUpConfirmHandler {
	return &SignUpConfirmHandler{
		confirmSignUpSvc: confirmSignUpSvc,
	}
}

func (h *SignUpConfirmHandler) Handle(eCtx echo.Context) error {
	payload := oas.ConfirmSignUp{}

	if err := eCtx.Bind(&payload); err != nil {
		return err
	}

	if err := validateSignUpConfirm(&payload); err != nil {
		return errors.NewValidationError(err)
	}

	ctx := eCtx.Request().Context()

	if err := h.confirmSignUpSvc.Confirm(ctx, authDomain.ConfirmSignUpInput{
		Token: payload.Token,
		Code:  payload.Code,
	}); err != nil {
		return mapError(err)
	}

	return eCtx.NoContent(http.StatusNoContent)
}
