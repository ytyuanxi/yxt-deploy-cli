package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
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
		initCluster()
	},
}

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Display cluster information",
	Long:  "Display cluster information",
	Run: func(cmd *cobra.Command, args []string) {
		infoOfCluster()
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

type ServerType string

const (
	MasterServer   ServerType = "master"
	MetaNodeServer ServerType = "metanode"
	DataNodeServer ServerType = "datanode"
)

type Status string

const (
	Running Status = "running"
	Stopped Status = "stopped"
)

type Service struct {
	ServerType    ServerType
	ContainerName string
	NodeIP        string
	Status        Status
}

func printTable(services []Service) {
	fmt.Println("Server Type  | Container Name | Node IP         | Status")
	fmt.Println("-----------------------------------------------------")
	for _, service := range services {
		fmt.Printf("%-12s | %-14s | %-13s | %s\n", service.ServerType, service.ContainerName, service.NodeIP, service.Status)
	}
}

func infoOfCluster() error {
	//
	config, err := readConfig()
	if err != nil {
		log.Fatal(err)
	}

	servers := []Service{}

	for id, node := range config.DeployHostsList.Master.Hosts {
		server := Service{}
		server.NodeIP = node
		server.ServerType = MasterServer
		server.Status = Stopped
		server.ContainerName = "master" + strconv.Itoa(id+1)
		servers = append(servers, server)
		ps, err := psContainerOnNode(RemoteUser, node)
		if err != nil {
			return err
		}
		containerArray := strings.Split(ps, " ")
		for _, container := range containerArray {
			if container == "master"+strconv.Itoa(id+1) {
				server.Status = Running
				break
			}
		}

	}

	for id, node := range config.DeployHostsList.MetaNode.Hosts {
		server := Service{}
		server.NodeIP = node
		server.ServerType = MetaNodeServer
		server.Status = Stopped
		server.ContainerName = "metanode" + strconv.Itoa(id+1)
		servers = append(servers, server)
		ps, err := psContainerOnNode(RemoteUser, node)
		if err != nil {
			return err
		}
		containerArray := strings.Split(ps, " ")
		for _, container := range containerArray {
			if container == "metanode"+strconv.Itoa(id+1) {
				server.Status = Running
				break
			}
		}

	}

	for id, node := range config.DeployHostsList.DataNode {
		server := Service{}
		server.NodeIP = node.Hosts
		server.ServerType = DataNodeServer
		server.Status = Stopped
		server.ContainerName = "datanode" + strconv.Itoa(id+1)
		servers = append(servers, server)
		ps, err := psContainerOnNode(RemoteUser, node.Hosts)
		if err != nil {
			return err
		}
		containerArray := strings.Split(ps, " ")
		for _, container := range containerArray {
			if container == "datanode"+strconv.Itoa(id+1) {
				server.Status = Running
				break
			}
		}

	}
	printTable(servers)
	return nil
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
