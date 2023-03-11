package service

import (
	"KillerFeature/ServerSide/internal"
	cconn "KillerFeature/ServerSide/internal/client_conn"
	"KillerFeature/ServerSide/internal/models"
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrorCreateConnWrap = "creating conn error"
)

type Service struct {
	sshBuilder cconn.CCBuilder
}

func NewService(sshBuilder cconn.CCBuilder) internal.Usecase {
	return &Service{
		sshBuilder: sshBuilder,
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
	// execLib := GetOSCommandLib
	fmt.Println(cc.Exec("TODO"))
	//
	// TODO: Close() когда задеплоится
	cc.Close()
	return nil
}
