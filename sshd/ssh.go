package sshd

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
	mcSshPath := os.Getenv("MC_SSH_PATH")

	privateKeyFile, err := os.ReadFile(mcSshPath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(privateKeyFile)
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
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return nil, err

	}
	return output, nil
}
