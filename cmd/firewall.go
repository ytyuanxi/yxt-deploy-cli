package cmd

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

func openRemotePortFirewall(hostname, username string, privateKeyPath string, port int) error {
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostname, 22), config)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	command := fmt.Sprintf("sudo firewall-cmd --zone=public --add-port=%d/tcp --permanent", port)
	err = session.Run(command)
	if err != nil {
		return err
	}

	return nil
}

func checkPortStatus(nodeUser, node string, port string) (string, error) {
	cmd := exec.Command("ssh", nodeUser+"@"+node, "firewall-cmd --list-all | grep "+port)
	fmt.Println(cmd)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Port %s %s is closed", node, port), err
	}
	//fmt.Println(string(output))
	return fmt.Sprintf("Port %s is open", port), nil
}
