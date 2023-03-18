package usecase

import (
	"errors"
	"fmt"

	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib/ubuntu"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
)

type DeployAppUsecase struct {
	sshBuilder cconn.CCBuilder
	tm         *taskmanager.Manager
}

func NewDeployAppUsecase(sshBuilder cconn.CCBuilder, tm *taskmanager.Manager) ucase.DeployAppUsecase {
	return &DeployAppUsecase{
		sshBuilder: sshBuilder,
		tm:         tm,
	}
}

func (u *DeployAppUsecase) DeployApp(creds *models.SshCreds) (uint64, error) {
	taskId, err := u.tm.AddTask(DeployAppProcessTask, creds.Addr, taskmanager.AuthData{Login: creds.Login, Password: creds.Password})
	if err != nil {
		return uint64(taskId), errors.Join(ucase.ErrAddingToTaskManager, err)
	}
	return uint64(taskId), nil
}

func DeployAppProcessTask(cc cconn.ClientConn) error {
	deployCommands := getDeployCommands(ubuntu.Ubuntu2204CommandLib{})
	if deployCommands == nil {
		return ucase.ErrorUnsupportedOS
	}

	for _, command := range deployCommands {

		fmt.Println("=========command========")
		fmt.Println(command)

		output, err := cc.Exec(command.String())

		fmt.Println(err)
		fmt.Println(string(output))
		fmt.Println("=========end========")

		if err != nil {
			switch {
			case errors.Is(err, cconn.ErrExitStatus):
				{
					return errors.Join(ucase.ErrExecuteDeployInstructions, err)
				}
			case errors.Is(err, cconn.ErrExitStatusMissing):
				{
					return errors.Join(ucase.ErrMissingStatusDeployInstructions, err)
				}
			case errors.Is(err, cconn.ErrOpenChannel):
				{
					return errors.Join(ucase.ErrCreateSession, err)
				}
			default:
				{
					return errors.Join(ucase.ErrUnknown, err)
				}
			}
		}
	}

	return nil
}
