package auth

import (
	"github.com/labstack/echo/v4"

	aos "gitlab.com/apart-deals/openapi/go/api"
)

type AuthHandler struct {
}

func (h *AuthHandler) SignUp(eCtx echo.Context) error {
	_ = aos.SignUp{}

	return nil
}

func (h *AuthHandler) SignIn(eCtx echo.Context) error {
	return nil
}

func (h *AuthHandler) SignOut(eCtx echo.Context) error {
	return nil
}
