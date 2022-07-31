package server

import (
	"apart-deal-api/pkg/api/handlers/auth"

	"github.com/labstack/echo/v4"
)

func CreateRouter(
	e *echo.Echo,
	authHandler *auth.AuthHandler,
) *echo.Router {
	r := echo.NewRouter(e)

	r.Add("GET", "/auth/sign-up", authHandler.SignUp)
	r.Add("GET", "/auth/sign-in", authHandler.SignIn)
	r.Add("GET", "/auth/sign-out", authHandler.SignOut)

	return r
}
