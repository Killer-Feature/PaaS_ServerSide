package handler

import (
	"net/http"
	"net/netip"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
)

var (
	HttpErrorBindingParams = "Error binding request params."
	HttpErrorValidateIP    = "Your ip address is invalid. Please provide global IPv4 address."
	HttpErrorValidatePort  = "Your port is invalid. Please provide non-zero port."
	HttpErrorValidateAddr  = "Your address is invalid. Please provide global IPv4 address and non-zero port."
	HttpErrInternal        = "Internal server error."
)

type DeployAppHandler struct {
	logger *zap.Logger
	u      ucase.DeployAppUsecase
}

func NewDeployAppHandler(logger *zap.Logger, u ucase.DeployAppUsecase) *DeployAppHandler {
	return &DeployAppHandler{logger: logger, u: u}
}

func (h *DeployAppHandler) DeployApp(c echo.Context) error {
	var req models.SshDeployAppReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorBindingParams)
	}

	ip, err := netip.ParseAddr(req.IP)
	if err != nil || !ip.Is4() || ip.IsLoopback() || ip.IsPrivate() {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorValidateIP)
	}
	if req.Port == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorValidatePort)
	}

	ipPort := netip.AddrPortFrom(ip, req.Port)
	if !ipPort.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrorValidateAddr)
	}

	err = h.u.DeployApp(&models.SshCreds{
		Addr:     ipPort,
		Password: req.Password,
		Login:    req.Login,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, HttpErrInternal)
	}

	return c.NoContent(http.StatusOK)
}
