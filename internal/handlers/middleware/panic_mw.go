package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (mw *CommonMiddleware) PanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if err := recover(); err != nil {
				requestId := GetRequestIdFromCtx(ctx)
				mw.logger.RequestError(requestId, "panic recovered: "+fmt.Sprint(err))
				mw.logger.Access(requestId, ctx.Request().Method, ctx.Request().RemoteAddr, ctx.Request().URL.Path, time.Duration(0))
				_ = ctx.JSON(http.StatusInternalServerError, struct {
					Error string `json:"error"`
				}{Error: "internal server error"})
				if err != nil {
					return
				}
			}
		}()
		return next(ctx)
	}
}
