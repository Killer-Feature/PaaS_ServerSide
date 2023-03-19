package client_conn

import (
	"errors"
	"net/netip"
)

type CCBuilder interface {
	CreateCC(ip netip.AddrPort, login, password string) (ClientConn, error)
}

type ClientConn interface {
	Exec(command string) ([]byte, error)
	Close() error
}

var (
	ErrOperation               = errors.New("operation error")
	ErrOpenChannel             = errors.New("target server rejects an OpenChannel request")
	ErrExitStatusMissing       = errors.New("remote server did not send an exit status")
	ErrExitStatus              = errors.New("the command completes unsuccessfully or is interrupted by a signal")
	ErrConnectionAlreadyClosed = errors.New("connection already closed")
	ErrUnknown                 = errors.New("")
)
