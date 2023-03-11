package handlers

import (
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"KillerFeature/ServerSide/internal"
	models "KillerFeature/ServerSide/internal/models"
)

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

func (h *Handler) Register(s *echo.Echo) {
	// Register http handlers

	s.POST("/deploy-app", h.DeployApp)
}

func (h *Handler) DeployApp(c echo.Context) error {
	var req models.SshDeployAppReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorBindingParams)
	}
	// TODO delete
	fmt.Println(req)
	fmt.Println(h.u.DeployApp(&models.SshCreds{
		IP:       req.IP,
		Password: req.Password,
		User:     req.User,
	}))
	return c.HTML(http.StatusOK, "hello")
}
