package handler

import (
	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers/middleware"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/netip"
	"time"
)

const (
	READ_BUFSIZE              = 1024
	WRITE_BUFSIZE             = 1024
	DEPLOY_PROGRESS_CHAN_SIZE = 1024
	READ_CREDS_MSG_TIMEOUT    = 60
)

var (
	HttpErrorBindingParams = "Переданы невалидные данные."
	HttpErrorValidateIP    = "Передан невалидный IP-адрес. Пожалуйста, укажите глобальный IPv4 адрес."
	HttpErrorValidatePort  = "Передан невалидный порт. Пожалуйста, укажите ненулевой порт."
	HttpErrorValidateAddr  = "Передан невалидный IP-адрес или порт. Пожалуйста, укажите глобальный IPv4 адрес и ненулевой порт."
	HttpErrInternal        = "Внутренняя ошибка сервера."
	HttpErrUpgradeToWS     = "Ошибка соединения с сервером"
)

var (
	errAddDeployTaskToTaskManager = "error adding deploy-task to task manager"
	errSetReadWSTimeout           = "error setting the read deadline on the underlying network connection"
	errCloseWSConn                = "error closing the underlying network connection without sending or waiting for a close message"
	errUpgradingToWS              = "error upgrading the HTTP server connection to the WebSocket protocol"
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
	// почекать куку
	// если кука есть -- чекаем стоит ли задача в очереди
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.RequestError(reqId, errUpgradingToWS+err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, HttpErrUpgradeToWS)
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			h.logger.RequestError(reqId, errCloseWSConn+err.Error())
		}
	}(ws)

	var req models.SshDeployAppReq

	err = ws.SetReadDeadline(time.Now().Add(READ_CREDS_MSG_TIMEOUT * time.Second))
	if err != nil {
		h.logger.RequestError(reqId, errSetReadWSTimeout+err.Error())
		_ = ws.WriteJSON(models.Error{Error: HttpErrUpgradeToWS, Code: http.StatusBadRequest})
		return nil
	}
	err = ws.ReadJSON(&req)
	if err != nil {
		_ = ws.WriteJSON(models.Error{Error: HttpErrorBindingParams, Code: http.StatusBadRequest})
		return nil
	}

	ip, err := netip.ParseAddr(req.IP)
	if err != nil || !ip.Is4() || ip.IsLoopback() || ip.IsPrivate() {
		_ = ws.WriteJSON(models.Error{Error: HttpErrorValidateIP, Code: http.StatusBadRequest})
		return nil
	}
	if req.Port == 0 {
		_ = ws.WriteJSON(models.Error{Error: HttpErrorValidatePort, Code: http.StatusBadRequest})
		return nil
	}

	ipPort := netip.AddrPortFrom(ip, req.Port)
	if !ipPort.IsValid() {
		_ = ws.WriteJSON(models.Error{Error: HttpErrorValidateAddr, Code: http.StatusBadRequest})
		return nil
	}

	progressChan := make(chan models.TaskProgressMsg, DEPLOY_PROGRESS_CHAN_SIZE)
	_, err = h.u.DeployApp(&models.SshCreds{
		Addr:     ipPort,
		Password: req.Password,
		Login:    req.Login,
	}, progressChan)

	if err != nil {
		_ = ws.WriteJSON(models.Error{Error: HttpErrInternal, Code: http.StatusInternalServerError})
		h.logger.RequestError(reqId, errAddDeployTaskToTaskManager+": "+err.Error())
		return nil
	}

	h.deployApp(reqId, ws, progressChan)
	return nil
}

func (h *DeployAppHandler) deployApp(reqId uint64, client *websocket.Conn, progressChan chan models.TaskProgressMsg) {
	for msg := range progressChan {
		err := client.WriteJSON(msg)
		if err != nil {
			h.logger.RequestError(reqId, err.Error())
		}
	}
}
