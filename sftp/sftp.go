package sftp

import (
	"fmt"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"path"
	"time"
)


func Connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))
	// set ssh connect config
	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

func Put(sftpClient *sftp.Client, srcFilePath, dstDirPath string) error {

	// open src file
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	defer srcFile.Close()

	// create dst file in remote host
	dstFileName := path.Base(srcFilePath)
	dstFile, err := sftpClient.Create(path.Join(dstDirPath, dstFileName))
	if err != nil {
		logrus.Fatal(err)
		return err
	}
	defer dstFile.Close()

	// copy from src file to dst file
	buf := make([]byte, 1024)
	for ; ; {
		n,_ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	logrus.Infof("copy file to remote server(%v) finished!", &sftpClient)

	return nil
}
