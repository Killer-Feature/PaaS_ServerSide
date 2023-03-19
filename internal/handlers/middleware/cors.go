package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func GetCorsConfig(allowOrigins []string, maxAge int) middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXCSRFToken},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS, echo.PUT},
		ExposeHeaders:    []string{echo.HeaderXCSRFToken},
		MaxAge:           maxAge,
	}
}
