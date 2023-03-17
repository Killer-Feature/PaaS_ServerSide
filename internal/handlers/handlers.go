package handlers

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"go.uber.org/zap"

	ucase "KillerFeature/ServerSide/internal"
	models "KillerFeature/ServerSide/internal/models"
)

//go:embed dist
var ui embed.FS

var (
	HttpErrorBindingParams                = "Error binding request params"
	HttpErrDeployInstruction              = "Error executing deploy huggin application instructions. Deploy instruction exited with not 0 status"
	HttpErrDeployInstructionMissingStatus = "Error executing deploy huggin application instructions.deploy instructions did not send exit status."
	HttpErrUnsupportedOS                  = "Unsupported operating system installed to your server"
	HttpErrorCreateSession                = "Your server rejects create new session request to execute deploy huggin app instructions. (A session is a remote execution of a program request)"
	HttpErrorCreateCConn                  = "Error creating ssh connection to your server"
	HttpErrInternal                       = "Internal server error"
)

type Handler struct {
	logger *zap.Logger
	u      ucase.Usecase
}

func NewHandler(logger *zap.Logger, u ucase.Usecase) *Handler {
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
	log, err := h.u.DeployApp(&models.SshCreds{
		IP:       req.IP,
		Port:     req.Port,
		Password: req.Password,
		User:     req.User,
	})

	if err != nil {
		switch {
		case errors.Is(err, ucase.ErrExecuteDeployInstructions):
			{
				return c.JSON(http.StatusBadRequest, models.SshDeployAppErrorResp{Log: string(log), Error: HttpErrDeployInstruction})
			}
		case errors.Is(err, ucase.ErrMissingStatusDeployInstructions):
			{
				return c.JSON(http.StatusBadRequest, models.SshDeployAppErrorResp{Log: string(log), Error: HttpErrDeployInstructionMissingStatus})
			}
		case errors.Is(err, ucase.ErrorUnsupportedOS):
			{
				return c.JSON(http.StatusBadRequest, models.SshDeployAppErrorResp{Log: string(log), Error: HttpErrUnsupportedOS})
			}
		case errors.Is(err, ucase.ErrCreateSession):
			{
				return c.JSON(http.StatusBadRequest, models.SshDeployAppErrorResp{Log: string(log), Error: HttpErrorCreateSession})
			}
		case errors.Is(err, ucase.ErrCreateClientConnection):
			{
				return c.JSON(http.StatusBadRequest, models.SshDeployAppErrorResp{Log: string(log), Error: HttpErrorCreateCConn})
			}
		default:
			{
				return c.JSON(http.StatusInternalServerError, models.SshDeployAppErrorResp{Log: string(log), Error: HttpErrInternal})
			}
		}
	}

	return c.JSON(http.StatusOK, models.SshDeployAppResp{Log: string(log)})
}

