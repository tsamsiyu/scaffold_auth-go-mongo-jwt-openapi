package auth

import (
	"github.com/labstack/echo/v4"
)

type RouteGroup *echo.Group

func RegisterSignUpRoute(g RouteGroup, signUpHandler *SignUpHandler) {
	v := *g
	v.POST("/auth/sign-up", signUpHandler.Handle)
}

func RegisterSignUpConfirmRoute(g RouteGroup, signUpConfirmHandler *SignUpConfirmHandler) {
	v := *g
	v.POST("/auth/sign-up-confirm", signUpConfirmHandler.Handle)
}
