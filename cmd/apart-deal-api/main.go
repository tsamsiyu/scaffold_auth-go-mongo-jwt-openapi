package main

import (
	"fmt"
	"os"

	"apart-deal-api/cmd/apart-deal-api/app"
	"apart-deal-api/pkg/logging"
)

func main() {
	logger, err := logging.FromEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err := app.NewConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err = app.Run(cfg, logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("API stopped")
}
