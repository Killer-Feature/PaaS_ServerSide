package handler

import (
	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers/middleware"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/netip"
)

var (
	HttpErrorBindingParams = "Переданы невалидные данные."
	HttpErrorValidateIP    = "Передан невалидный IP-адрес. Пожалуйста, укажите глобальный IPv4 адрес."
	HttpErrorValidatePort  = "Передан невалидный порт. Пожалуйста, укажите ненулевой порт."
	HttpErrorValidateAddr  = "Передан невалидный IP-адрес или порт. Пожалуйста, укажите глобальный IPv4 адрес и ненулевой порт."
	HttpErrInternal        = "Внутренняя ошибка сервера."
)

var (
	errAddDeployTaskToTaskManager = "error adding deploy-task to task manager"
)

type DeployAppHandler struct {
	logger *servlog.ServLogger
	u      ucase.DeployAppUsecase
}

func NewDeployAppHandler(logger *servlog.ServLogger, u ucase.DeployAppUsecase) *DeployAppHandler {
	return &DeployAppHandler{logger: logger, u: u}
}

func (h *DeployAppHandler) DeployApp(c echo.Context) error {
	reqId := middleware.GetRequestIdFromCtx(c)
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

	//TODO: set deploy-task-id cookie
	_, err = h.u.DeployApp(&models.SshCreds{
		Addr:     ipPort,
		Password: req.Password,
		Login:    req.Login,
	})

	if err != nil {
		h.logger.RequestError(reqId, errAddDeployTaskToTaskManager+": "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, HttpErrInternal)
	}

	return c.NoContent(http.StatusOK)
}
