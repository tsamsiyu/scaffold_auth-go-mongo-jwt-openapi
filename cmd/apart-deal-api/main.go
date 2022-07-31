package main

import (
	"fmt"
	"os"

	"apart-deal-api/cmd/apart-deal-api/app"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed while creating logger: ", err)
		os.Exit(1)
	}

	if err = app.Run(&app.Config{
		Port: 8080,
	}, logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info("API stopped")
}
