package handlers

import (
	"embed"
	deployappdelivery "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app/delivery"
	"github.com/labstack/echo/v4"
	"io/fs"
	"net/http"
)

//go:embed dist
var ui embed.FS

func Register(s *echo.Echo, deployAppHandler deployappdelivery.DeployAppHandler) error {
	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		return err
	}

	s.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
	s.POST("/deploy-app", deployAppHandler.DeployApp)
	return nil
}
