package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster manager",
	Long:  `This command will manager the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.UsageString())
	},
}

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "init the cluster from config.yaml",
	Long:  "init the cluster from config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		// Read Configuration File
		config, err := readConfig()
		if err != nil {
			log.Fatal(err)
		}

		hosts := []string{}
		hosts = append(hosts, config.DeployHostsList.Master.Hosts...)
		hosts = append(hosts, config.DeployHostsList.MetaNode.Hosts...)
		for i := 0; i < len(config.DeployHostsList.DataNode); i++ {
			hosts = append(hosts, config.DeployHostsList.DataNode[i].Hosts)
		}

		// Obtain the IP address of the current host
		currentNode, err := getCurrentIP()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("The IP address of the current host:", currentNode)

		// Establish a secure connection from the current node to other nodes
		for _, node := range hosts {
			if node == currentNode || node == "" {
				continue
			}
			err := establishSSHConnectionWithoutPassword(currentNode, "root", node)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("Password free connection establishment completed")

		for _, node := range hosts {
			// Check if Docker is installed and installed
			if node == currentNode || node == "" {
				continue
			}
			checkAndInstallDocker("root", node)

			// Check if the Docker service is started and started
			err = checkAndStartDockerService("root", node)
			if err != nil {
				log.Printf("Failed to start Docker service on node% s:% v", node, err)
			} else {
				log.Printf("The docker for node %s is ready", node)
			}

			// Pull Mirror
			err = pullImageOnNode("root", node, config.Global.ContainerImage)
			if err != nil {
				log.Printf("Failed to pull mirror% s on node% s:% v", node, config.Global.ContainerImage, err)
			} else {
				log.Printf("Successfully pulled mirror % s on node % s", config.Global.ContainerImage, node)
			}
		}

		log.Println("*******Cluster environment initialization completed******")
	},
}

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Display cluster information",
	Long:  "Display cluster information",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var clearCommand = &cobra.Command{
	Use:   "clear",
	Short: "Clear cluster files and information",
	Long:  "Clear cluster files and information",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	ClusterCmd.AddCommand(initCommand)
	ClusterCmd.AddCommand(infoCommand)
	ClusterCmd.AddCommand(clearCommand)
}

// 检查是否已存在私钥和公钥文件并生成SSH密钥对
// Check if the private and public key files already exist and generate an SSH key pair
func generateSSHKey() error {
	// Check if the private and public key files already exist
	privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	//publicKeyPath := privateKeyPath + ".pub"
	if _, err := os.Stat(privateKeyPath); err == nil {
		return fmt.Errorf("SSH key already exists")
	}

	// Generate SSH key pairs
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-N", "", "-f", privateKeyPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to generate SSH key: %v", err)
	}

	log.Printf("SSH key generated successfully.\n")

	return nil
}

// Establishing a secure connection
func establishSSHConnectionWithoutPassword(sourceNode, targetNodeUser, targetNode string) error {

	// Check if the private and public key files already exist and generate an SSH key pairs
	generateSSHKey()

	// Check if it is possible to connect to the target node without a password
	cmd := exec.Command("ssh", "-o", "BatchMode=yes", "-o", "ConnectTimeout=5", targetNodeUser+"@"+targetNode, "echo", "connection successful")
	err := cmd.Run()
	if err != nil {
		// Copy public key to remote host
		privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
		publicKeyPath := privateKeyPath + ".pub"
		cmd = exec.Command("ssh-copy-id", "-i", publicKeyPath, targetNodeUser, "@", targetNode)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to establish passwordless SSH connection with %s@%s: %v", targetNodeUser, targetNode, err)
		}
	}
	log.Printf("Passwordless SSH connection is established with %s@%s.\n", targetNodeUser, targetNode)

	return nil
}

// 检查Docker是否安装并安装
// Check if Docker is installed and installed
func checkAndInstallDocker(nodeUser, node string) error {
	// Check if Docker is installed
	cmd := exec.Command("ssh", nodeUser+"@"+node, "docker --version")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "Docker version") {
		//log.Println("Docker installed")
		return nil
	}

	// Docker not installed, installing Docker
	cmd = exec.Command("ssh", nodeUser+"@"+node, "yum", "install", "docker", "-y")

	// Set output to standard output and standard error output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute Command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to install Docker on node %s", node)
	}

	return nil

}

// 检查Docker服务是否启动并启动
// Check if the Docker service is started and started
func checkAndStartDockerService(nodeUser, node string) error {
	// Check Docker Service Status
	cmd := exec.Command("ssh", nodeUser+"@"+node, "systemctl is-active docker.service")
	output, err := cmd.Output()

	if err == nil && strings.TrimSpace(string(output)) == "active" {
		// Docker service started
		//log.Println("docker already start")
		return nil
	}

	// Docker service not started, starting Docker service
	cmd = exec.Command("ssh", nodeUser+"@"+node, "systemctl start docker")

	// Set output to standard output and standard error output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to start Docker service on node %s", node)
	}
	log.Println("docker start")

	return nil
}

// 从配置文件中拉取镜像
// Pull image from configuration file
func pullImageOnNode(nodeUser, node, imageName string) error {
	// Remote execution of commands to pull images
	cmd := exec.Command("ssh", nodeUser+"@"+node, "docker pull "+imageName)

	//Set output to standard output and standard error output
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to pull mirror %s on node %s", imageName, node)
	}

	return nil
}

// 获取当前主机的IP地址
// Obtain the IP address of the current host
func getCurrentIP() (string, error) {
	// Get Host Name
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	//fmt.Println(hostname)

	// Obtain the IP address of the host
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}
	//fmt.Println(addrs)

	// Select IPv4 address
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			//fmt.Println(ipv4.String())
			return ipv4.String(), nil

		}
	}

	return "", fmt.Errorf("IPv4 address not found")
}
