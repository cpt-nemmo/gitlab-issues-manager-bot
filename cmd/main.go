package main

import (
	"gitlab-issues-manager/internal/di"
	"go.uber.org/zap"
	"log"
)

func main() {
	di := di.DI{}

	err := di.Init()
	if err != nil {
		log.Fatal("failed to initialize DI", zap.Error(err))
	}
	di.StartBot()
}
