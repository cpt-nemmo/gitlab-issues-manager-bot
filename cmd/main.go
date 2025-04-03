package main

import (
	"go.uber.org/zap"
	"log"
	"test/internal/di"
)

func main() {
	di := di.DI{}

	err := di.Init()
	if err != nil {
		log.Fatal("failed to initialize DI", zap.Error(err))
	}
	di.StartBot()
}
