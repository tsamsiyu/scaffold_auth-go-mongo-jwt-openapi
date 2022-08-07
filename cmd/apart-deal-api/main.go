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

	envs, err := app.NewEnvs()
	if err != nil {
		logger.Fatal(err.Error())
	}

	if err = app.Run(envs, logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("API stopped")
}
