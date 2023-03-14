package client_conn

import (
	"KillerFeature/ServerSide/internal/models"
	"errors"
)

type CCBuilder interface {
	CreateCC(creds *models.SshCreds) (ClientConn, error)
}

type ClientConn interface {
	Exec(command string) ([]byte, error)
	Close() error
}

var (
	ErrOperation               = errors.New("operation error")
	ErrOpenChannel             = errors.New("target server rejects an OpenChannel request")
	ErrExecCommand             = errors.New("error executing command")
	ErrExitStatusMissing       = errors.New("remote server does not send an exit status")
	ErrExitStatus              = errors.New("the command completes unsuccessfully or is interrupted by a signal")
	ErrConnectionAlreadyClosed = errors.New("connection already closed")
	ErrSome                    = errors.New("")
)
