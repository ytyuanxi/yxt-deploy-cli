package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

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
		initCluster()
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

func initCluster() {
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
		if node == "" {
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

		err = transferFileToRemote("bin", config.Global.DataDir, "root", node)
		if err != nil {
			log.Println(err)
		}
		err = transferFileToRemote("script", config.Global.DataDir, "root", node)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("*******Cluster environment initialization completed******")
}
