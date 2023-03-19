package ssh

import (
	"errors"
	"net"
	"net/netip"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"

	cc "github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn"
)

const (
	SSH_CCONN_TIMEOUT = 60
)

type SSHBuilder struct {
}

type SSH struct {
	C *ssh.Client
}

func NewSSHBuilder() cc.CCBuilder {
	return &SSHBuilder{}
}

func (b *SSHBuilder) CreateCC(addr netip.AddrPort, login, password string) (cc.ClientConn, error) {
	clientConfig := ssh.ClientConfig{
		User: login,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyAlgorithms: []string{
			ssh.CertAlgoRSASHA256v01,
			ssh.CertAlgoRSAv01,
			ssh.CertAlgoRSASHA512v01,
			ssh.KeyAlgoRSASHA256,
			ssh.KeyAlgoRSASHA512,
			ssh.KeyAlgoRSA,
			ssh.CertAlgoECDSA256v01,
			ssh.KeyAlgoED25519,
			ssh.CertAlgoED25519v01,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * SSH_CCONN_TIMEOUT,
	}

	client, err := ssh.Dial("tcp", addr.String(), &clientConfig)
	if err != nil {
		var opErrTarget *net.OpError
		if errors.As(err, &opErrTarget) {
			return nil, errors.Join(cc.ErrOperation, err)
		}
		return nil, errors.Join(cc.ErrUnknown, err)
	}

	return &SSH{
		C: client,
	}, nil
}

func (s *SSH) Exec(command string) ([]byte, error) {
	session, err := s.C.NewSession()
	if err != nil {
		var openChannelErrTarget *ssh.OpenChannelError
		if errors.As(err, &openChannelErrTarget) {
			return nil, errors.Join(cc.ErrOpenChannel, err)
		}
		return nil, errors.Join(cc.ErrUnknown, err)
	}

	output, err := session.CombinedOutput(command)

	if err != nil {
		var exitMissingErrTarget *ssh.ExitMissingError
		if errors.As(err, &exitMissingErrTarget) {
			return nil, errors.Join(cc.ErrExitStatusMissing, err)
		}

		var exitErrTarget *ssh.ExitError
		if errors.As(err, &exitErrTarget) {
			return nil, errors.Join(cc.ErrExitStatus, err)
		}

		return nil, errors.Join(cc.ErrUnknown, err)
	}

	return output, err
}

func (s *SSH) Close() error {
	err := s.C.Close()
	if err != nil {
		var opErrTarget *net.OpError
		if errors.As(err, &opErrTarget) {
			return errors.Join(cc.ErrOperation, err)
		}
		alreadyClosedErrTarget := syscall.EINVAL
		if errors.As(err, &alreadyClosedErrTarget) {
			return errors.Join(cc.ErrConnectionAlreadyClosed, err)
		}
		return errors.Join(cc.ErrUnknown, err)
	}
	return nil
}
