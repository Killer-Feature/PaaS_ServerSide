package middleware

import (
	"github.com/labstack/echo/v4"
	"time"
)

const LoggerCtxKey = "logger"

func (mw *CommonMiddleware) AccessLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reqId := GetRequestIdFromCtx(ctx)
		ctx.Set(LoggerCtxKey, mw.logger)
		start := time.Now()
		result := next(ctx)
		mw.logger.Access(reqId, ctx.Request().Method, ctx.Request().RemoteAddr, ctx.Request().URL.Path, time.Since(start))
		return result
	}
}
