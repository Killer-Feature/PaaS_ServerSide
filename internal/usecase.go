package internal

import (
	"KillerFeature/ServerSide/internal/models"
	"errors"
)

type Usecase interface {
	DeployApp(creds *models.SshCreds) (string, error)
}

var (
	ErrCreateClientConnection          = errors.New("error creating client connection")
	ErrorUnsupportedOS                 = errors.New("unsupported operating system")
	ErrUnknown                         = errors.New("")
	ErrCreateSession                   = errors.New("target server rejects create new session request")
	ErrMissingStatusDeployInstructions = errors.New("deploy instructions did not send exit status")
	ErrExecuteDeployInstructions       = errors.New("deploy instruction exited with not 0 status")
)
