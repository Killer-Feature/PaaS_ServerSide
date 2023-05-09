package usecase

import (
	"bytes"
	"errors"
	ssh2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/key_value_storage"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	cl2 "github.com/Killer-Feature/PaaS_ServerSide/pkg/os_command_lib"
	"net/netip"

	ucase "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	cconn "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
)

type DeployAppUsecase struct {
	sshBuilder      cconn.CCBuilder
	tm              *taskmanager.Manager[netip.AddrPort]
	logger          *servlog.ServLogger
	progressStorage key_value_storage.KeyValueStorage[netip.Addr, models.TaskProgressMsg]
}

func NewDeployAppUsecase(sshBuilder cconn.CCBuilder, tm *taskmanager.Manager[netip.AddrPort], ps key_value_storage.KeyValueStorage[netip.Addr, models.TaskProgressMsg]) *DeployAppUsecase {
	return &DeployAppUsecase{
		sshBuilder:      sshBuilder,
		tm:              tm,
		progressStorage: ps,
	}
}

type progressRecorder struct {
	ip     netip.Addr
	taskId uint64
	c      *chan models.TaskProgressMsg
	s      key_value_storage.KeyValueStorage[netip.Addr, models.TaskProgressMsg]
	logger *servlog.ServLogger
}

func (p *progressRecorder) commitTaskProcess(status models.DeployAppStatus, percent uint8, log []byte, errStr string) {
	msg := models.TaskProgressMsg{
		Status:  status,
		Percent: percent,
		Log:     string(log),
		Error:   errStr,
		Chan:    p.c,
		TaskId:  p.taskId,
	}
	*p.c <- msg
	var err error
	if status == models.STATUS_SUCCESS {
		err = p.s.DeleteByKey(&p.ip)
	}
	err = p.s.Set(&p.ip, &msg)
	if err != nil {
		p.logger.TaskError(p.taskId, ucase.ErrUpdateProgress.Error()+": "+err.Error())
	}
}

func pushToLog(log []byte, command []byte, output []byte) []byte {
	return bytes.Join([][]byte{log, append([]byte("$ "), command...), output}, []byte("\n"))
}

func (u *DeployAppUsecase) DeployApp(creds *models.SshCreds) (uint64, error) {
	progressChan := make(chan models.TaskProgressMsg, ucase.DEPLOY_PROGRESS_CHAN_SIZE)
	taskId, err := u.tm.AddTask(u.DeployAppProcessTask(creds, &progressChan), creds.Addr)
	if err != nil {
		err = errors.Join(ucase.ErrAddingToTaskManager, err)
		u.logger.TaskError(uint64(taskId), err.Error())
		return uint64(taskId), err
	}

	progressRec := progressRecorder{ip: creds.Addr.Addr(), taskId: uint64(taskId), logger: u.logger, c: &progressChan, s: u.progressStorage}
	progressRec.commitTaskProcess(models.STATUS_IN_QUEUE, 0, nil, "")
	return uint64(taskId), nil
}

func (u *DeployAppUsecase) DeployAppProcessTask(creds *models.SshCreds, progressChan *chan models.TaskProgressMsg) func(taskId taskmanager.ID) error {
	return func(taskId taskmanager.ID) error {
		progressRec := progressRecorder{ip: creds.Addr.Addr(), taskId: uint64(taskId), logger: u.logger, c: progressChan, s: u.progressStorage}
		var percent uint8 = 1
		defer close(*progressChan)

		progressRec.commitTaskProcess(models.STATUS_START, percent, nil, "")

		sshBuilder := ssh2.NewSSHBuilder()
		cc, err := sshBuilder.CreateCC(creds.Addr, creds.Login, creds.Password)
		percent = 5
		if err != nil {
			err = errors.Join(ucase.ErrCreateCC, err)
			progressRec.commitTaskProcess(models.STATUS_CONN_ERR, percent, nil, err.Error())
			return err
		}
		defer func(cc cconn.ClientConn) {
			err := cc.Close()
			if err != nil {
				u.logger.TaskError(uint64(taskId), ucase.ErrCloseCC.Error()+": "+err.Error())
			}
		}(cc)

		progressRec.commitTaskProcess(models.STATUS_IN_PROCESS, percent, nil, "")

		osRelease, log, err := getOSRelease(cc)
		percent = 10
		if err != nil {
			err = errors.Join(ucase.ErrUnsupportedOS, err)
			progressRec.commitTaskProcess(models.STATUS_ERROR, percent, log, err.Error())
			return err
		}

		progressRec.commitTaskProcess(models.STATUS_IN_PROCESS, percent, log, "")

		percent = 11
		deployCommands := getDeployCommands(osRelease, creds.Login, creds.Password)
		if deployCommands == nil {
			progressRec.commitTaskProcess(models.STATUS_ERROR, percent, log, ucase.ErrUnsupportedOS.Error())
			return ucase.ErrUnsupportedOS
		}

		percentStep := (uint8(99) - percent) / uint8(len(deployCommands))

		for _, command := range deployCommands {
			output, err := cc.Exec(command.String())
			log = pushToLog(log, []byte(command.Command), output)
			if err != nil {
				switch {
				case errors.Is(err, cconn.ErrExitStatus):
					{
						err = errors.Join(ucase.ErrExecuteDeployInstructions, err)
						progressRec.commitTaskProcess(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				case errors.Is(err, cconn.ErrExitStatusMissing):
					{
						err = errors.Join(ucase.ErrMissingStatusDeployInstructions, err)
						progressRec.commitTaskProcess(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				case errors.Is(err, cconn.ErrOpenChannel):
					{
						err = errors.Join(ucase.ErrCreateSession, err)
						progressRec.commitTaskProcess(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				default:
					{
						err = errors.Join(ucase.ErrUnknown, err)
						progressRec.commitTaskProcess(models.STATUS_ERROR, percent, log, err.Error())
						return err
					}
				}
			}

			percent += percentStep
			progressRec.commitTaskProcess(models.STATUS_IN_PROCESS, percent, log, "")
		}

		progressRec.commitTaskProcess(models.STATUS_SUCCESS, 100, log, "")

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
func (u *DeployAppUsecase) ProgressInfo(ip *netip.Addr) (*models.TaskProgressMsg, error) {
	progressInfo, err := u.progressStorage.GetByKey(ip)
	if err != nil {
		if errors.Is(err, key_value_storage.ErrNoSuchElem) {
			return nil, errors.Join(err, ucase.ErrSuchIPISNotProcessing)
		}
		err = errors.Join(err, ucase.ErrUnknown)
		u.logger.TaskError(uint64(progressInfo.TaskId), err.Error())
		return nil, err
	}
	return progressInfo, nil
}
