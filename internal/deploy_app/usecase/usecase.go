package usecase

import (
	"errors"
	"fmt"
	ssh2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	cl2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"
	"net/netip"

	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
)

type DeployAppUsecase struct {
	sshBuilder cconn.CCBuilder
	tm         *taskmanager.Manager[netip.AddrPort]
	logger     *servlog.ServLogger
}

func NewDeployAppUsecase(sshBuilder cconn.CCBuilder, tm *taskmanager.Manager[netip.AddrPort]) ucase.DeployAppUsecase {
	return &DeployAppUsecase{
		sshBuilder: sshBuilder,
		tm:         tm,
	}
}

func (u *DeployAppUsecase) DeployApp(creds *models.SshCreds) (uint64, error) {
	taskId, err := u.tm.AddTask(u.DeployAppProcessTask(creds), creds.Addr)
	if err != nil {
		return uint64(taskId), errors.Join(ucase.ErrAddingToTaskManager, err)
	}
	return uint64(taskId), nil
}

func (u *DeployAppUsecase) DeployAppProcessTask(creds *models.SshCreds) func(taskId taskmanager.ID) error {
	return func(taskId taskmanager.ID) error {
		creds := creds
		sshBuilder := ssh2.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(creds.Addr, creds.Login, creds.Password)
		if err != nil {
			return errors.Join(ucase.ErrCreateCC, err)
		}
		defer func(cc cconn.ClientConn) {
			err := cc.Close()
			if err != nil {
				u.logger.TaskError(uint64(taskId), ucase.ErrCloseCC.Error()+": "+err.Error())
			}
		}(cc)

		osRelease, err := getOSRelease(cc)
		if err != nil {
			return errors.Join(ucase.ErrorUnsupportedOS, err)
		}

		deployCommands := getDeployCommands(osRelease)
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
}

func getOSRelease(cc cconn.ClientConn) (cl2.OSRelease, error) {
	osReleaseCommand := cl2.GetOSRelease()
	output, err := cc.Exec(string(osReleaseCommand.Command))

	if err != nil {
		switch {
		case errors.Is(err, cconn.ErrExitStatus):
			{
				return cl2.UnknownOS, errors.Join(ucase.ErrExecuteDeployInstructions, err)
			}
		case errors.Is(err, cconn.ErrExitStatusMissing):
			{
				return cl2.UnknownOS, errors.Join(ucase.ErrMissingStatusDeployInstructions, err)
			}
		case errors.Is(err, cconn.ErrOpenChannel):
			{
				return cl2.UnknownOS, errors.Join(ucase.ErrCreateSession, err)
			}
		default:
			{
				return cl2.UnknownOS, errors.Join(ucase.ErrUnknown, err)
			}
		}
	}
	osRelease := cl2.UnknownOS
	_ = osReleaseCommand.Parser(output, &osRelease)
	return osRelease, nil
}
