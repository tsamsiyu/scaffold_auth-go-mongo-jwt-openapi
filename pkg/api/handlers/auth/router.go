package auth

import (
	"github.com/labstack/echo/v4"
)

type RouteGroup *echo.Group

func RegisterSignUpRoute(g RouteGroup, signUpHandler *SignUpHandler) {
	v := *g
	v.POST("/sign-up", signUpHandler.Handle)
}

func RegisterSignUpConfirmRoute(g RouteGroup, signUpConfirmHandler *SignUpConfirmHandler) {
	v := *g
	v.POST("/sign-up-confirm", signUpConfirmHandler.Handle)
}

func RegisterSignInRoute(g RouteGroup, signInHandler *SignInHandler) {
	v := *g
	v.POST("/sign-in", signInHandler.Handle)
}
