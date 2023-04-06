package handler

import (
	"errors"
	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers/middleware"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"net/netip"
	"time"
)

const (
	READ_BUFSIZE           = 1024
	WRITE_BUFSIZE          = 1024
	DEPLOY_APP_COOKIE_NAME = "deploy_ip"
)

var (
	HttpErrBindingParams         = "Переданы невалидные данные."
	HttpErrValidateIP            = "Передан невалидный IP-адрес. Пожалуйста, укажите глобальный IPv4 адрес."
	HttpErrValidatePort          = "Передан невалидный порт. Пожалуйста, укажите ненулевой порт."
	HttpErrValidateAddr          = "Передан невалидный IP-адрес или порт. Пожалуйста, укажите глобальный IPv4 адрес и ненулевой порт."
	HttpErrInternal              = "Внутренняя ошибка сервера."
	HttpErrSuchIPISNotProcessing = "Не найдена информация по процессу деплоя huginn на данном ip адресе"
)

var (
	errAddDeployTaskToTaskManager = "error adding deploy-task to task manager"
	errCloseWSConn                = "error closing the underlying network connection without sending or waiting for a close message"
	errUpgradingToWS              = "error upgrading the HTTP server connection to the WebSocket protocol"
	errGetProgressInfo            = "error getting progress info from storage by cookie"
)

type DeployAppHandler struct {
	logger *servlog.ServLogger
	u      ucase.DeployAppUsecase
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  READ_BUFSIZE,
	WriteBufferSize: WRITE_BUFSIZE,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	//	TODO: add CheckOrigin func(r *http.Request) bool field to avoid CSRF
}

func NewDeployAppHandler(logger *servlog.ServLogger, u ucase.DeployAppUsecase) *DeployAppHandler {
	return &DeployAppHandler{logger: logger, u: u}
}

func (h *DeployAppHandler) DeployApp(c echo.Context) error {
	reqId := middleware.GetRequestIdFromCtx(c)

	var req models.SshDeployAppReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrBindingParams)
	}

	ip, err := netip.ParseAddr(req.IP)
	if err != nil || !ip.Is4() || ip.IsLoopback() || ip.IsPrivate() {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrValidateIP)
	}
	if req.Port == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrValidatePort)
	}

	ipPort := netip.AddrPortFrom(ip, req.Port)
	if !ipPort.IsValid() {
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrValidateAddr)
	}

	_, err = h.u.DeployApp(&models.SshCreds{
		Addr:     ipPort,
		Password: req.Password,
		Login:    req.Login,
	})

	if err != nil {
		h.logger.RequestError(reqId, errAddDeployTaskToTaskManager+": "+err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, HttpErrInternal)
	}

	host, _, _ := net.SplitHostPort(c.Request().Host)
	c.SetCookie(&http.Cookie{
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Name:     DEPLOY_APP_COOKIE_NAME,
		Value:    ip.String(),
		Domain:   host,
		Path:     "/",
	})

	return c.NoContent(http.StatusOK)
}

func (h *DeployAppHandler) Deploying(c echo.Context) error {
	reqId := middleware.GetRequestIdFromCtx(c)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.RequestError(reqId, errUpgradingToWS+err.Error())
		return nil
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			h.logger.RequestError(reqId, errCloseWSConn+err.Error())
		}
	}(ws)

	ipCookie, err := c.Cookie(DEPLOY_APP_COOKIE_NAME)
	if err != nil {
		return nil
	}

	ip, err := netip.ParseAddr(ipCookie.Value)
	if err != nil || !ip.Is4() || ip.IsLoopback() || ip.IsPrivate() {
		_ = ws.WriteJSON(models.TaskProgressMsg{Error: HttpErrValidateIP, Percent: 0, Status: models.STATUS_ERROR})
		return nil
	}

	progressInfo, err := h.u.ProgressInfo(&ip)

	if err != nil {
		if errors.Is(err, ucase.ErrSuchIPISNotProcessing) {
			h.logger.RequestError(reqId, errGetProgressInfo+": "+err.Error())
			_ = ws.WriteJSON(models.TaskProgressMsg{Error: HttpErrSuchIPISNotProcessing, Percent: 0, Status: models.STATUS_ERROR})
			return nil
		}
		h.logger.RequestError(reqId, errGetProgressInfo+": "+err.Error())
		_ = ws.WriteJSON(models.TaskProgressMsg{Error: HttpErrInternal, Percent: 0, Status: models.STATUS_ERROR})
		return nil
	}

	h.deploying(reqId, ws, progressInfo)
	return nil
}

func (h *DeployAppHandler) deploying(reqId uint64, client *websocket.Conn, progressInfo *models.TaskProgressMsg) {
	err := client.WriteJSON(progressInfo)
	if err != nil {
		h.logger.RequestError(reqId, err.Error())
	}

	var msg models.TaskProgressMsg
	for msg = range *(progressInfo.Chan) {
		if msg.Percent > progressInfo.Percent {
			err = client.WriteJSON(msg)
			if err != nil {
				h.logger.RequestError(reqId, err.Error())
			}
		}
	}
}
