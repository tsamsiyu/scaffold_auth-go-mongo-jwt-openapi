package main

import (
	"apart-deal-api/cmd/apart-deal-api/app"
	"apart-deal-api/dependencies"
)

func main() {
	logger := dependencies.LoggerFromEnv()

	if err := app.Run(logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("API stopped")
}
