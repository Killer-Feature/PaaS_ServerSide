package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"KillerFeature/ServerSide/internal"
	models "KillerFeature/ServerSide/internal/models"
)

//go:embed dist
var ui embed.FS

var (
	HttpErrorBindingParams = "Error binding request params"
)

type Handler struct {
	logger *zap.Logger
	u      internal.Usecase
}

func NewHandler(logger *zap.Logger, u internal.Usecase) *Handler {
	return &Handler{logger: logger, u: u}
}

type Response struct {
	StatusCode int    `json:"status_code"`
	Data       string `json:"data"`
}

func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers

	s.POST("/deploy-app", h.DeployApp)

	fsys, err := fs.Sub(ui, "dist")
	if err != nil {
		h.logger.Fatal("fs creating error", zap.Error(err))
	}

	s.GET("/*", echo.WrapHandler(http.FileServer(http.FS(fsys))))
}

func (h *Handler) DeployApp(c echo.Context) error {
	var req models.SshDeployAppReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorBindingParams)
	}
	// TODO: валидация
	err := h.u.DeployApp(&models.SshCreds{
		IP:       req.IP,
		Port:     req.Port,
		Password: req.Password,
		User:     req.User,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			StatusCode: http.StatusInternalServerError,
			Data:       err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}
