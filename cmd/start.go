package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var ip string
var allStart bool
var datanodeDisk string

const ClusterName = "cubeFS"

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start service",
	Long:  `This command will start services.`,
	Run: func(cmd *cobra.Command, args []string) {
		if allStart {
			err := startALLMaster()
			if err != nil {
				log.Println(err)
			}

		} else {
			fmt.Println(cmd.UsageString())
		}

	},
}

var startMasterCommand = &cobra.Command{
	Use:   "master",
	Short: "start master",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("start master in ", ip)
		} else {
			fmt.Println("start all master services from config.yaml")
		}
	},
}

var startMetanodeCommand = &cobra.Command{
	Use:   "metanode",
	Short: "start metanode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("start metanode in ", ip)
		} else {
			fmt.Println("start all metanode services from config.yaml")
		}

	},
}

var startDatanodeCommand = &cobra.Command{
	Use:   "datanode",
	Short: "start datanode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("start datanode in ", ip)
			if !cmd.Flags().Changed("disk") {
				fmt.Println("must have disk argument")
				os.Exit(1)
			}
			fmt.Println("disk:", datanodeDisk)
		} else {
			fmt.Println("start all datanode services from config.yaml")
		}
	},
}

func init() {
	StartCmd.AddCommand(startMasterCommand)
	StartCmd.AddCommand(startMetanodeCommand)
	StartCmd.AddCommand(startDatanodeCommand)
	StartCmd.Flags().BoolVarP(&allStart, "all", "a", false, "start all services")
	StartCmd.PersistentFlags().StringVarP(&ip, "ip", "", "", "specify an IP address to start services")
	startDatanodeCommand.Flags().StringVarP(&datanodeDisk, "disk", "d", "", "specify the disk where datanode mount")
}

func startALLMaster() error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	masterListenPort := config.Master.Config.Listen
	masterProfPort := config.Master.Config.Prof
	masterDataDir := config.Master.Config.DataDir

	peers := ""
	////获取master的peers "1:192.168.0.11:17010,2:192.168.0.12:17010,3:192.168.0.13:17010"
	for id, node := range config.DeployHostsList.Master.Hosts {
		if id != len(config.DeployHostsList.Master.Hosts)-1 {
			peers = peers + strconv.Itoa(id+1) + ":" + node + ":" + config.Master.Config.Listen + ","
		} else {
			peers = peers + strconv.Itoa(id+1) + ":" + node + ":" + config.Master.Config.Listen
		}

	}
	fmt.Println(peers)

	//对每个节点：scp相应的文件到该节点，在该节点启动相应的容器
	for id, node := range config.DeployHostsList.Master.Hosts {
		//检查服务所对应端口的防火墙是否开放
		listenStatus, err := checkPortStatus("root", node, masterListenPort)
		log.Println(listenStatus)
		if err != nil {
			//开放该端口
			privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
			portNum, _ := strconv.Atoi(masterListenPort)
			err = openRemotePortFirewall(node, "root", privateKeyPath, portNum)
			if err != nil {
				return fmt.Errorf("firewall opened for prot failed")
			}
			fmt.Printf("Firewall opened for port %d successfully.\n", portNum)
		}

		profStatus, err := checkPortStatus("root", node, masterProfPort)
		fmt.Println(profStatus)
		if err != nil {
			//开放该端口
			privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
			portNum, _ := strconv.Atoi(masterListenPort)
			err = openRemotePortFirewall(node, "root", privateKeyPath, portNum)
			if err != nil {
				return fmt.Errorf("firewall opened for prot failed")
			}
			fmt.Printf("Firewall opened for port %d successfully.\n", portNum)
			return err
		}

		//传输bin文件目标节点

		err = scpFile(config.Global.BinDir, masterDataDir, node, "22")
		if err != nil {
			return err
		}
		fmt.Println("bin file transferred successfully.")

		//解析yaml文件为master.json文件
		err = writeMaster(ClusterName, strconv.Itoa(id+1), node, masterListenPort, masterProfPort, peers)
		if err != nil {
			return err
		}

		//并将该文件传输到目标节点

		err = scpFile("master.json", "master.json", node, "22")
		if err != nil {
			return err
		}
		//在该节点启动容器
		err = startMasterContainer("master"+strconv.Itoa(id+1), node)
		if err != nil {
			return err
		}
	}

	fmt.Println("start all master services from config.yaml")

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

func scpFile(localPath string, remotePath string, hostname string, port string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "put %s\n", remotePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, file)
	if err != nil {
		return err
	}

	return nil
}

func startMasterContainer(containerName, node string) error {
	cmd := exec.Command("docker", "start", containerName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

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

	// Replace "PORT" with the actual port number you want to open
	command := fmt.Sprintf("sudo firewall-cmd --zone=public --add-port=%d/tcp --permanent", port)
	err = session.Run(command)
	if err != nil {
		return err
	}

	return nil
}
