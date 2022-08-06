package auth

import (
	"github.com/labstack/echo/v4"

	authDomain "apart-deal-api/pkg/domain/auth"
	oas "gitlab.com/apart-deals/openapi/go/api"
)

type AuthHandler struct {
	signUpSvc        *authDomain.SignUpService
	confirmSignUpSvc *authDomain.ConfirmSignUpService
}

func NewAuthHandler(signUpSvc *authDomain.SignUpService, confirmSignUpSvc *authDomain.ConfirmSignUpService) *AuthHandler {
	return &AuthHandler{
		signUpSvc:        signUpSvc,
		confirmSignUpSvc: confirmSignUpSvc,
	}
}

func (h *AuthHandler) SignUp(eCtx echo.Context) error {
	payload := oas.SignUp{}

	if err := eCtx.Bind(&payload); err != nil {
		return err
	}

	ctx := eCtx.Request().Context()

	if err := h.signUpSvc.SignUp(ctx, authDomain.SignUpInput{
		Email:    payload.Email,
		Name:     payload.Name,
		Password: payload.Password,
	}); err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) ConfirmSignUp(eCtx echo.Context) error {
	payload := oas.ConfirmSignUp{}

	if err := eCtx.Bind(&payload); err != nil {
		return err
	}

	ctx := eCtx.Request().Context()

	if err := h.confirmSignUpSvc.Confirm(ctx, authDomain.ConfirmSignUpInput{
		Token: payload.Token,
		Code:  payload.Code,
	}); err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) SignIn(eCtx echo.Context) error {
	return nil
}

func (h *AuthHandler) SignOut(eCtx echo.Context) error {
	return nil
}
