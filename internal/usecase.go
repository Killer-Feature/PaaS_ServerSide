package internal

import (
	"KillerFeature/ServerSide/internal/models"
	"errors"
)

type Usecase interface {
	DeployApp(creds *models.SshCreds) (*models.SshDeployAppServiceResp, error)
}

var (
	ErrCreateClientConnaction = errors.New("error creating client connection")
	ErrorUnsupportedOS        = errors.New("unsupported operating system")
	ErrSome                   = errors.New("")
)
