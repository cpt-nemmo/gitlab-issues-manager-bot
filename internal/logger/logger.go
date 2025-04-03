package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func Enter(method string) *zap.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Sampling = nil
	l, _ := cfg.Build()

	l.Info(fmt.Sprintf("enter %s", method))

	return l
}

func Exit(log *zap.Logger, method string) {
	log.Info(fmt.Sprintf("exit %s", method))
}
