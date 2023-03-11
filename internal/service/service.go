package service

import (
	"KillerFeature/ServerSide/internal"
	cconn "KillerFeature/ServerSide/internal/client_conn"
	"KillerFeature/ServerSide/internal/models"
	"KillerFeature/ServerSide/pkg/taskmanager"
)

var (
	ErrorCreateConnWrap = "creating conn error"
)

type Service struct {
	sshBuilder cconn.CCBuilder
	tm         *taskmanager.Manager
}

func NewService(sshBuilder cconn.CCBuilder, tm *taskmanager.Manager) internal.Usecase {
	return &Service{
		sshBuilder: sshBuilder,
		tm:         tm,
	}
}

func (s *Service) DeployApp(creds *models.SshCreds) error {
	// TODO: валидация IP -- белый список
	// TODO: task ID

	_, err := s.tm.AddTask(creds.IP+":"+creds.Port, creds.User, creds.Password)
	if err != nil {
		return err
	}

	//cc, err := s.sshBuilder.CreateCC(creds)
	//if err != nil {
	//	return errors.Wrap(err, ErrorCreateConnWrap)
	//}
	// TODO: GetOSCommandLib возвращает структуру с командами для конкретной ОС, вид ОС можно узнать через SSH
	// execLib := GetOSCommandLib
	//fmt.Println(cc.Exec("TODO"))

	//
	// TODO: Close() когда задеплоится
	//cc.Close()
	return nil
}
