package service

import (
	"KillerFeature/ServerSide/internal"
	cconn "KillerFeature/ServerSide/internal/client_conn"
	"KillerFeature/ServerSide/internal/models"
	command_lib "KillerFeature/ServerSide/pkg/os_command_lib"
	ubuntu2204_commands "KillerFeature/ServerSide/pkg/os_command_lib/ubuntu2204"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrorCreateConnWrap   = "creating conn error"
	ErrorDeployingAppWrap = "deploying app error"
)

var (
	ErrorUnsupportedOS = errors.New("unsupported operating system")
)

const (
	HUGGIN_DIR        = "huggin"
	HUGGIN_BINARY_URL = "https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_22.04"
)

type Service struct {
	sshBuilder cconn.CCBuilder
}

func NewService(sshBuilder cconn.CCBuilder) internal.Usecase {
	return &Service{
		sshBuilder: sshBuilder,
	}
}

type os string

const (
	ubuntu2204 os = "UBUNTU 22.04"
)

func getDeployCommands(os os) ([]command_lib.Command, error) {
	switch os {
	case ubuntu2204:
		return []command_lib.Command{
			// "rm -rf Huginn",
			// "curl \"127.0.0.1:8090\"",
			// "kill $(lsof -t -i:8090)",
			"if lsof -i :8090; then kill $(lsof -i :8090); fi",
			// ubuntu2204_commands.Mkdir.WithArgs("-p", HUGGIN_DIR),
			// ubuntu2204_commands.Cd.WithArgs(HUGGIN_DIR),
			ubuntu2204_commands.Wget.WithArgs(HUGGIN_BINARY_URL),
			ubuntu2204_commands.Chmod.WithArgs("777", "PaaS_22.04"), // TODO concat filename,
			"nohup ./PaaS_22.04 &",                                  //TODO
			// "curl \"127.0.0.1:8090\"",
		}, nil
	default:
		{
			return nil, ErrorUnsupportedOS
		}
	}
}

func (s *Service) DeployApp(creds *models.SshCreds) error {
	// TODO: валидация IP -- белый список
	// TODO: task ID

	cc, err := s.sshBuilder.CreateCC(creds)
	if err != nil {
		return errors.Wrap(err, ErrorCreateConnWrap)
	}
	// TODO: GetOSCommandLib возвращает структуру с командами для конкретной ОС, вид ОС можно узнать через SSH
	// execLib := GetOSCommandLib("")

	deployCommands, err := getDeployCommands(ubuntu2204)
	if err != nil {
		return err
	}

	for _, command := range deployCommands {
		output, err := cc.Exec(command.String())
		fmt.Println(string(output), err)
		if err != nil {
			return errors.Wrap(err, ErrorDeployingAppWrap)
		}
	}
	// mkdir -p huginn
	// cd huginn
	// wget https://github.com/Killer-Feature/PaaS_ClientSide/releases/download/v0.0.1/PaaS_22.04
	// chmod 777 PaaS_22.04*
	// ./PaaS_22.04

	// TODO: Close() когда задеплоится
	cc.Close()
	return nil
}
