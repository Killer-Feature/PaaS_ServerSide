package ssh

import (
	cc "KillerFeature/ServerSide/internal/client_conn"
	models "KillerFeature/ServerSide/internal/models"
	"net"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
)

const (
	SSH_CCONN_TIMEOUT = 60
)

var (
	ErrorCreateCCEmptyCreds = errors.New("empty user or password passed")
	ErrorCreateCCEmptyAddr  = errors.New("empty address or port passed")
	ErrorDialCConn          = errors.New("starting ssh client connection error")
	ErrorOpenNewSession     = errors.New("opening new session error")
)

type SSHBuilder struct {
}

type SSH struct {
	// C ssh.Client
	S *ssh.Session
}

func NewSSHBuilder() cc.CCBuilder {
	return &SSHBuilder{}
}

func getHostKeyCallback() ssh.HostKeyCallback {
	return ssh.InsecureIgnoreHostKey()
}

func (b *SSHBuilder) CreateCC(creds *models.SshCreds) (cc.ClientConn, error) {
	if creds.User == "" || creds.Password == "" {
		return nil, ErrorCreateCCEmptyCreds
	}
	if creds.IP == "" || creds.Port == "" {
		return nil, ErrorCreateCCEmptyAddr
	}

	hostKeyCallback := getHostKeyCallback()

	var authSlice []ssh.AuthMethod
	authSlice = append(authSlice, ssh.Password(creds.Password))

	clientConfig := ssh.ClientConfig{
		User: creds.User,
		Auth: authSlice,
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
		HostKeyCallback: hostKeyCallback,
		Timeout:         time.Second * SSH_CCONN_TIMEOUT,
	}
	dial, err := ssh.Dial("tcp", net.JoinHostPort(creds.IP, creds.Port), &clientConfig)
	if err != nil {
		return nil, ErrorDialCConn
	}

	session, err := dial.NewSession()
	if err != nil {
		return nil, ErrorOpenNewSession
	}

	return &SSH{
		S: session,
	}, nil
}

func (s *SSH) Exec(comand string) ([]byte, error) {
	// TODO: сделать асинхронный деплой приложения
	output, err := s.S.Output("ls")
	return output, err
}

func (s *SSH) Close() error {
	return s.S.Close()
}
