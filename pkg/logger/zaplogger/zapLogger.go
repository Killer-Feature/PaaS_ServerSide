package zaplogger

import (
	"go.uber.org/zap"
)

func NewZapLogger(zapCfg *zap.Config) (*zap.SugaredLogger, error) {
	loggerUnsugared, err := zap.Config(*zapCfg).Build()
	if err != nil {
		return nil, err
	}
	return loggerUnsugared.Sugar(), nil
}
