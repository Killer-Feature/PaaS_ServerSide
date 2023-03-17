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
	tm         *taskmanager.Manager
}

func NewService(sshBuilder cconn.CCBuilder, tm *taskmanager.Manager) internal.Usecase {
	return &Service{
		sshBuilder: sshBuilder,
		tm:         tm,
	}
}

func pushToLog(log []byte, command []byte, output []byte) []byte {
	return bytes.Join([][]byte{log, append([]byte("$ "), command...), output}, []byte("\n"))
}

func (s *Service) DeployApp(creds *models.SshCreds) (string, error) {
	_, err := s.tm.AddTask(creds.IP+":"+creds.Port, creds.User, creds.Password)
	return err
}
