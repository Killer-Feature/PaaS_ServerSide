package handlers

import (
	"embed"
	deployappdelivery "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app/delivery"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers/middleware"
	echo "github.com/labstack/echo/v4"
	"io/fs"
	"net/http"
)

//go:embed dist
var ui embed.FS

func Register(s *echo.Echo, deployAppHandler deployappdelivery.DeployAppHandler, mw *middleware.CommonMiddleware) error {
	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		return err
	}

	mwChain := []echo.MiddlewareFunc{
		mw.RequestIdMiddleware,
		mw.AccessLogMiddleware,
		mw.PanicMiddleware,
	}

	s.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))), mwChain...)
	s.POST("/deploy-app", deployAppHandler.DeployApp, mwChain...)
	return nil
}
