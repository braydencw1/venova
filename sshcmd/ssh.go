package sshcmd

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnectToDev() (*ssh.Client, error) {
	mcHost := os.Getenv("MC_HOST")
	mcPort := os.Getenv("MC_PORT")
	mcUser := os.Getenv("MC_USER")
	mcSshKey := os.Getenv("MC_SSH_KEY")

	signer, err := ssh.ParsePrivateKey([]byte(mcSshKey))
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: mcUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", mcHost, mcPort), config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func RunCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	defer func() { _ = session.Close() }()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return nil, err

	}
	return output, nil
}
