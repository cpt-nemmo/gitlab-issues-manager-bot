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

func Exit(log *zap.Logger, method string, err error) {
	if err != nil {
		log.Error(fmt.Sprintf("exit %s", method), zap.Error(err))
		return
	}

	log.Info(fmt.Sprintf("exit %s", method))
}
