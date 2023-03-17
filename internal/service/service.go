package service

import (
	"KillerFeature/ServerSide/internal"
	ucase "KillerFeature/ServerSide/internal"
	cconn "KillerFeature/ServerSide/internal/client_conn"
	"KillerFeature/ServerSide/internal/models"
	"KillerFeature/ServerSide/pkg/os_command_lib/ubuntu"
	"bytes"
	"errors"
	"fmt"
)

const (
	INITIAL_LOG_LEN = 1024
)

type Service struct {
	sshBuilder cconn.CCBuilder
	//	TODO: add logger & log errors
}

func NewService(sshBuilder cconn.CCBuilder) internal.Usecase {
	return &Service{
		sshBuilder: sshBuilder,
	}
}

func pushToLog(log []byte, command []byte, output []byte) []byte {
	return bytes.Join([][]byte{log, append([]byte("$ "), command...), output}, []byte("\n"))
}

func (s *Service) DeployApp(creds *models.SshCreds) (string, error) {
	cc, err := s.sshBuilder.CreateCC(creds)
	if err != nil {
		if errors.Is(err, cconn.ErrOperation) {
			return "", errors.Join(ucase.ErrCreateClientConnection, err)
		}
		return "", errors.Join(ucase.ErrUnknown, err)
	}
	defer func(cc cconn.ClientConn) {
		err := cc.Close()
		if err != nil {
			//	TODO: log error
		}
	}(cc)

	// TODO: GetOSCommandLib возвращает структуру с командами для конкретной ОС, вид ОС можно узнать через SSH
	// execLib := GetOSCommandLib("")

	deployCommands := getDeployCommands(ubuntu.Ubuntu2204CommandLib{})
	if deployCommands == nil {
		return "", ucase.ErrorUnsupportedOS
	}

	log := make([]byte, 0, INITIAL_LOG_LEN)
	for _, command := range deployCommands {

		fmt.Print("=========command========")
		fmt.Print(command)
		fmt.Println("=========end========")

		output, err := cc.Exec(command.String())
		log = pushToLog(log, []byte(command), output)

		fmt.Println(err)
		fmt.Println(string(output))
		if err != nil {
			switch {
			case errors.Is(err, cconn.ErrExitStatus):
				{
					return string(log), errors.Join(ucase.ErrExecuteDeployInstructions, err)
				}
			case errors.Is(err, cconn.ErrExitStatusMissing):
				{
					return string(log), errors.Join(ucase.ErrMissingStatusDeployInstructions, err)
				}
			case errors.Is(err, cconn.ErrOpenChannel):
				{
					return string(log), errors.Join(ucase.ErrCreateSession, err)
				}
			default:
				{
					return string(log), errors.Join(ucase.ErrUnknown, err)
				}
			}
		}
	}

	return string(log), nil
}
