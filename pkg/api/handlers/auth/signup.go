package auth

import (
	"net/http"

	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/labstack/echo/v4"

	apiErr "apart-deal-api/pkg/api/aspects/errors"
	authDomain "apart-deal-api/pkg/domain/auth"

	validation "github.com/go-ozzo/ozzo-validation"
	oas "gitlab.com/apart-deals/openapi/go/api"
)

func validateSignUp(payload *oas.SignUp) error {
	return validation.ValidateStruct(
		payload,
		validation.Field(&payload.Name, validation.Required, validation.Length(2, 50)),
		validation.Field(&payload.Email, validation.Required, is.Email, validation.Length(3, 50)),
		validation.Field(&payload.Password, validation.Required, validation.Length(4, 10)),
	)
}

type SignUpHandler struct {
	signUpSvc *authDomain.SignUpService
}

func NewSignUpHandler(signUpSvc *authDomain.SignUpService) *SignUpHandler {
	return &SignUpHandler{
		signUpSvc: signUpSvc,
	}
}

func (h *SignUpHandler) Handle(eCtx echo.Context) error {
	payload := oas.SignUp{}

	if err := eCtx.Bind(&payload); err != nil {
		return err
	}

	if err := validateSignUp(&payload); err != nil {
		return apiErr.NewMultipleValidationInputError(err)
	}

	ctx := eCtx.Request().Context()

	output, err := h.signUpSvc.SignUp(ctx, authDomain.SignUpInput{
		Email:    payload.Email,
		Name:     payload.Name,
		Password: payload.Password,
	})
	if err != nil {
		return mapError(err)
	}

	return eCtx.JSON(http.StatusOK, oas.SignUpResponse{
		Token: output.Token,
	})
}
