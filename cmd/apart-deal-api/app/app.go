package app

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

type Config struct {
	Port int
}

func Run(cfg *Config, logger *zap.Logger) error {
	e := echo.New()

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		return err
	}

	return nil
}
