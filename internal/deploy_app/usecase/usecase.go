package usecase

import (
	"bytes"
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

func newTaskProcessMsg(status models.DeployAppStatus, percent uint8, log []byte, err string) models.TaskProgressMsg {
	return models.TaskProgressMsg{
		Status:  status,
		Percent: percent,
		Log:     string(log),
		Error:   err,
	}
}

func pushToLog(log []byte, command []byte, output []byte) []byte {
	return bytes.Join([][]byte{log, append([]byte("$ "), command...), output}, []byte("\n"))
}

func (u *DeployAppUsecase) DeployApp(creds *models.SshCreds, progressChan chan models.TaskProgressMsg) (uint64, error) {
	taskId, err := u.tm.AddTask(u.DeployAppProcessTask(creds, progressChan), creds.Addr)
	if err != nil {
		return uint64(taskId), errors.Join(ucase.ErrAddingToTaskManager, err)
	}

	progressChan <- newTaskProcessMsg(models.STATUS_IN_QUEUE, 0, nil, "")

	return uint64(taskId), nil
}

func (u *DeployAppUsecase) DeployAppProcessTask(creds *models.SshCreds, progressChan chan models.TaskProgressMsg) func(taskId taskmanager.ID) error {
	return func(taskId taskmanager.ID) error {
		var percent uint8 = 0
		defer close(progressChan)

		percent = 1
		progressChan <- newTaskProcessMsg(models.STATUS_START, percent, nil, "")

		sshBuilder := ssh2.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(creds.Addr, creds.Login, creds.Password)
		percent = 5
		if err != nil {
			err = errors.Join(ucase.ErrCreateCC, err)
			progressChan <- newTaskProcessMsg(models.STATUS_CONN_ERR, percent, nil, errors.Join(ucase.ErrCreateCC, err).Error())
			return err
		}
		defer func(cc cconn.ClientConn) {
			err := cc.Close()
			if err != nil {
				u.logger.TaskError(uint64(taskId), ucase.ErrCloseCC.Error()+": "+err.Error())
			}
		}(cc)

		fmt.Println(creds)
		progressChan <- newTaskProcessMsg(models.STATUS_IN_PROCESS, percent, nil, "")

		osRelease, log, err := getOSRelease(cc)
		percent = 10
		if err != nil {
			err = errors.Join(ucase.ErrorUnsupportedOS, err)
			progressChan <- newTaskProcessMsg(models.STATUS_ERROR, percent, log, err.Error())
			return err
		}
		progressChan <- newTaskProcessMsg(models.STATUS_IN_PROCESS, percent, log, "")

		deployCommands := getDeployCommands(osRelease)
		if deployCommands == nil {
			progressChan <- newTaskProcessMsg(models.STATUS_ERROR, percent, log, ucase.ErrorUnsupportedOS.Error())
			return ucase.ErrorUnsupportedOS
		}

		percentStep := (uint8(99) - percent) / uint8(len(deployCommands))

		for _, command := range deployCommands {

			fmt.Println("=========command========")
			fmt.Println(command)

			output, err := cc.Exec(command.String())

			fmt.Println(err)
			fmt.Println(string(output))
			fmt.Println("=========end========")

			log = pushToLog(log, []byte(command.Command), output)

			if err != nil {
				switch {
				case errors.Is(err, cconn.ErrExitStatus):
					{
						err = errors.Join(ucase.ErrExecuteDeployInstructions, err)
						progressChan <- newTaskProcessMsg(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				case errors.Is(err, cconn.ErrExitStatusMissing):
					{
						err = errors.Join(ucase.ErrMissingStatusDeployInstructions, err)
						progressChan <- newTaskProcessMsg(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				case errors.Is(err, cconn.ErrOpenChannel):
					{
						err = errors.Join(ucase.ErrCreateSession, err)
						progressChan <- newTaskProcessMsg(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				default:
					{
						err = errors.Join(ucase.ErrUnknown, err)
						progressChan <- newTaskProcessMsg(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				}
			}

			percent += percentStep
			progressChan <- newTaskProcessMsg(models.STATUS_IN_PROCESS, percent, log, "")

		}

		progressChan <- newTaskProcessMsg(models.STATUS_SUCCESS, 100, log, "")
		return nil
	}
}

func getOSRelease(cc cconn.ClientConn) (cl2.OSRelease, []byte, error) {
	osReleaseCommand := cl2.GetOSRelease()
	output, err := cc.Exec(string(osReleaseCommand.Command))

	var log []byte
	log = pushToLog(log, []byte(osReleaseCommand.Command), output)

	if err != nil {
		switch {
		case errors.Is(err, cconn.ErrExitStatus):
			{
				return cl2.UnknownOS, log, errors.Join(ucase.ErrExecuteDeployInstructions, err)
			}
		case errors.Is(err, cconn.ErrExitStatusMissing):
			{
				return cl2.UnknownOS, log, errors.Join(ucase.ErrMissingStatusDeployInstructions, err)
			}
		case errors.Is(err, cconn.ErrOpenChannel):
			{
				return cl2.UnknownOS, log, errors.Join(ucase.ErrCreateSession, err)
			}
		default:
			{
				return cl2.UnknownOS, log, errors.Join(ucase.ErrUnknown, err)
			}
		}
	}
	osRelease := cl2.UnknownOS
	_ = osReleaseCommand.Parser(output, &osRelease)
	return osRelease, log, nil
}
