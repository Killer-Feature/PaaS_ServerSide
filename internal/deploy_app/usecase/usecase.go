package usecase

import (
	ucase "KillerFeature/ServerSide/internal/deploy_app"
	"KillerFeature/ServerSide/internal/models"
	cconn "KillerFeature/ServerSide/pkg/client_conn"
	"KillerFeature/ServerSide/pkg/os_command_lib/ubuntu"
	"KillerFeature/ServerSide/pkg/taskmanager"
	"errors"
	"fmt"
)

type DeployAppUsecase struct {
	sshBuilder cconn.CCBuilder
	//	TODO: add logger & log errors
	tm *taskmanager.Manager
}

func NewDeployAppUsecase(sshBuilder cconn.CCBuilder, tm *taskmanager.Manager) ucase.DeployAppUsecase {
	return &DeployAppUsecase{
		sshBuilder: sshBuilder,
		tm:         tm,
	}
}

func (s *DeployAppUsecase) DeployApp(creds *models.SshCreds) error {
	_, err := s.tm.AddTask(DeployAppProcessTask, creds.Addr, taskmanager.AuthData{Login: creds.Login, Password: creds.Password})
	if err != nil {
		return errors.Join(ucase.ErrAddingToTaskManager, err)
	}
	return nil
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
