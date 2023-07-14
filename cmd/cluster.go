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
		// 读取配置文件
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

		// 获取当前主机的IP地址
		currentNode, err := getCurrentIP()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("当前主机的IP地址:", ip)

		// 建立当前节点到其他节点的免密连接
		for _, node := range hosts {
			err := establishSSHConnection(currentNode, node)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("免密连接建立完成")

		for _, node := range hosts {
			// 检查Docker是否安装并安装
			err = checkAndInstallDocker(node)
			if err != nil {
				log.Printf("在节点 %s 上安装Docker失败: %v", node, err)
			} else {
				log.Printf("在节点 %s 上成功安装Docker", node)
			}

			// 检查Docker服务是否启动并启动
			err = checkAndStartDockerService(node)
			if err != nil {
				log.Printf("在节点 %s 上启动Docker服务失败: %v", node, err)
			} else {
				log.Printf("在节点 %s 上成功启动Docker服务", node)
			}

			// 拉取镜像
			err = pullImageOnNode(node, config.Global.ContainerImage)
			if err != nil {
				log.Printf("在节点 %s 上拉取镜像 %s 失败: %v", node, imageName, err)
			} else {
				log.Printf("在节点 %s 上成功拉取镜像 %s", node, imageName)
			}
		}

		fmt.Println("集群环境初始化完成")
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

// 建立免密连接
func establishSSHConnection(sourceNode, targetNode string) error {
	// 生成SSH密钥
	cmd := exec.Command("ssh-keygen", "-t", "rsa", "-N", "", "-f", "id_rsa")
	cmd.Dir = os.Getenv("HOME")
	err := cmd.Run()
	if err != nil {
		return err
	}

	// 将公钥复制到目标节点的authorized_keys文件中
	cmd = exec.Command("ssh-copy-id", targetNode)
	cmd.Dir = os.Getenv("HOME")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// 检查Docker是否安装
func checkAndInstallDocker(node string) error {
	// 检查Docker是否已安装
	cmd := exec.Command("ssh", node, "docker --version")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "Docker version") {
		// Docker已安装
		return nil
	}

	// Docker未安装，安装Docker
	cmd = exec.Command("ssh", node, "curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh")

	// 设置输出到标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("在节点 %s 上安装Docker失败: %v", node, err)
	}

	return nil

}

// 检查Docker服务是否启动
func checkAndStartDockerService(node string) error {
	// 检查Docker服务状态
	cmd := exec.Command("ssh", node, "systemctl is-active docker.service")
	output, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(output)) == "active" {
		// Docker服务已启动
		return nil
	}

	// Docker服务未启动，启动Docker服务
	cmd = exec.Command("ssh", node, "sudo systemctl start docker.service")

	// 设置输出到标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("在节点 %s 上启动Docker服务失败: %v", node, err)
	}

	return nil
}

// 从配置文件中拉取镜像
func pullImageOnNode(node, imageName string) error {
	// 远程执行拉取镜像的命令
	cmd := exec.Command("ssh", node, "docker pull "+imageName)

	// 设置输出到标准输出和标准错误输出
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 执行命令
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("在节点 %s 上拉取镜像 %s 失败: %v", node, imageName, err)
	}

	return nil
}

// 获取当前主机的IP地址
func getCurrentIP() (string, error) {
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// 获取主机的IP地址
	addrs, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	// 选择IPv4地址
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", fmt.Errorf("未找到IPv4地址")
}
