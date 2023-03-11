package client_conn

import (
	"net/netip"
)

type CCBuilder interface {
	// CreateCC() ClientConn
	CreateCC(creds *Creds) (ClientConn, error)
}

type Creds struct {
	IP       netip.AddrPort
	Login    string
	Password string
}

type ClientConn interface {
	Exec(command string) ([]byte, error)
	Close() error
}
