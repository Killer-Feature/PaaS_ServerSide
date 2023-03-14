package ssh

import (
	cc "KillerFeature/ServerSide/internal/client_conn"
	models "KillerFeature/ServerSide/internal/models"
	"net"
	"syscall"
	"time"

	"errors"

	"golang.org/x/crypto/ssh"
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

func (b *SSHBuilder) CreateCC(creds *models.SshCreds) (cc.ClientConn, error) {
	clientConfig := ssh.ClientConfig{
		User: creds.User,
		Auth: []ssh.AuthMethod{ssh.Password(creds.Password)},
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

	client, err := ssh.Dial("tcp", net.JoinHostPort(creds.IP, creds.Port), &clientConfig)
	if err != nil {
		var opErrTarget *net.OpError
		if errors.As(err, &opErrTarget) {
			return nil, errors.Join(cc.ErrOperation, err)
		}
		return nil, errors.Join(cc.ErrSome, err)
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
		return nil, errors.Join(cc.ErrSome, err)
	}

	output, err := session.CombinedOutput(command)

	if err != nil {
		var exitMissingErrTarget *ssh.ExitMissingError
		if errors.As(err, &exitMissingErrTarget) {
			return nil, errors.Join(cc.ErrExitStatusMissing, err)
		}

		var exitErrTarget *ssh.ExitMissingError
		if errors.As(err, &exitErrTarget) {
			return nil, errors.Join(cc.ErrExitStatus, err)
		}

		return nil, errors.Join(cc.ErrSome, err)
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
		return errors.Join(cc.ErrSome, err)
	}
	return nil
}
