package client_conn

import (
	"KillerFeature/ServerSide/internal/models"
)

type CCBuilder interface {
	// CreateCC() ClientConn
	CreateCC(creds *models.SshCreds) (ClientConn, error)
}

type ClientConn interface {
	Exec(command string) ([]byte, error)
	Close() error
}
