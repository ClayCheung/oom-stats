package ssh

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

func Connect(user string, password string, host string, port int) (*ssh.Session, error){
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	// get connect address
	addr = fmt.Sprintf("%s:%d", host, port)
	// set ssh connect config
	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	// connect to ssh
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	// create ssh session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil
}

func RunCmd(session *ssh.Session, cmd string) (string, error) {
	// run cmd list, return output list
	stdout, err := session.Output(cmd)
	if err != nil {
		logrus.Errorln(err)
	}
	return string(stdout), nil
}