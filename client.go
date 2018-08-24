package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Host     string
	Port     int
	User     string
	Password string
	KeyFile  string

	sshClient *ssh.Client
}

func NewSSHClient(host string, port int, user string, password string, keyfile string) *SSHClient {
	return &SSHClient{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		KeyFile:  keyfile,
	}
}

func (s *SSHClient) RunCommand(cmd string) (string, error) {
	ses, e := s.NewSession()
	if e != nil {
		e = fmt.Errorf("Unable to connect: %s", e.Error())
		return "", e
	}
	defer ses.Close()

	var out bytes.Buffer
	ses.Stdout = &out
	var err bytes.Buffer
	ses.Stderr = &err
	if e := ses.Run(cmd); e != nil {
		return "", errors.New(fmt.Sprintf("%s. %s", e.Error(), err.String()))
	}

	return out.String(), nil
}

func (s *SSHClient) NewSession() (*ssh.Session, error) {

	if s.sshClient == nil {
		c, err := s.connect()
		if err != nil {
			err = fmt.Errorf("Unable to connect: %s", err.Error())
			return nil, err
		}
		s.sshClient = c
	}

	ses, err := s.sshClient.NewSession()
	if err != nil {
		err = fmt.Errorf("Unable to start new session: %s", err.Error())
		return ses, err
	}

	return ses, err
}

func (s *SSHClient) connect() (*ssh.Client, error) {

	addr := net.JoinHostPort(s.Host, fmt.Sprintf("%d", s.Port))

	var auth []ssh.AuthMethod
	if s.Password != "" {
		auth = append(auth, ssh.Password(s.Password))
	}
	keyAuth := publicKeyFile(s.KeyFile)
	if keyAuth != nil {
		auth = append(auth, keyAuth)
	}

	cfg := &ssh.ClientConfig{
		User:            s.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            auth,
	}

	return ssh.Dial("tcp", addr, cfg)
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
