package deploy_app

import (
	"errors"

	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
)

type DeployAppUsecase interface {
	DeployApp(creds *models.SshCreds) (uint64, error)
}

var (
	ErrAddingToTaskManager             = errors.New("error adding deploy-task to task manager")
	ErrorUnsupportedOS                 = errors.New("unsupported operating system")
	ErrUnknown                         = errors.New("")
	ErrCreateSession                   = errors.New("target server rejects create new session request")
	ErrMissingStatusDeployInstructions = errors.New("deploy instructions did not send exit status")
	ErrExecuteDeployInstructions       = errors.New("deploy instruction exited with not 0 status")
	ErrCreateCC                        = errors.New("error creating new ssh connection")
	ErrCloseCC                         = errors.New("error closing ssh connection")
)
