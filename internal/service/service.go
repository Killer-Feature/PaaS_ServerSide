package service

import (
	"KillerFeature/ServerSide/internal"
	ucase "KillerFeature/ServerSide/internal"
	cconn "KillerFeature/ServerSide/internal/client_conn"
	"KillerFeature/ServerSide/internal/models"
	"errors"
	"fmt"
)

const (
	INITIAL_LOG_LEN = 1024
)

type Service struct {
	sshBuilder cconn.CCBuilder
}

func NewService(sshBuilder cconn.CCBuilder) internal.Usecase {
	return &Service{
		sshBuilder: sshBuilder,
	}
}

func pushToLog(log []byte, command []byte, output []byte) []byte {
	return fmt.Sprintf("%s\n$%s\n%s\n", log, command, output)
}

func (s *Service) DeployApp(creds *models.SshCreds) (*models.SshDeployAppServiceResp, error) {
	cc, err := s.sshBuilder.CreateCC(creds)
	if err != nil {
		if errors.Is(err, cconn.ErrOperation) {
			return nil, errors.Join(ucase.ErrCreateClientConnaction, err)
		}
		return nil, errors.Join(ucase.ErrSome, err)
	}
	defer cc.Close()

	// TODO: GetOSCommandLib возвращает структуру с командами для конкретной ОС, вид ОС можно узнать через SSH
	// execLib := GetOSCommandLib("")

	deployCommands, err := getDeployCommands(ubuntu2204)
	if err != nil {
		return nil, err
	}

	// если произошла ошибка в процессе выполнения команды -- выдавать весь вывод
	log := make([]byte, INITIAL_LOG_LEN)
	for _, command := range deployCommands {
		output, err := cc.Exec(command.String())
		log = pushToLog(log)
		fmt.Println("-------", string(command), "-------")
		fmt.Println(string(output), err)
		// if err != nil {
		// return errors.Wrap(err, ErrorDeployingAppWrap)
		// }
	}

	return nil, nil
}
