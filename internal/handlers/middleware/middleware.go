package middleware

import servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"

type CommonMiddleware struct {
	logger *servlog.ServLogger
}

func NewCommonMiddleware(logger *servlog.ServLogger) *CommonMiddleware {
	return &CommonMiddleware{
		logger: logger,
	}
}
