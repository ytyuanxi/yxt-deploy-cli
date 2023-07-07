package cmd

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Global          GlobalConfig          `yaml:"global"`
	Master          MasterConfig          `yaml:"master"`
	MetaNode        MetaNodeConfig        `yaml:"metanode"`
	DataNode        DataNodeConfig        `yaml:"datanode"`
	DeployHostsList DeployHostsListConfig `yaml:"deplopy_hosts_list"`
}

type GlobalConfig struct {
	SSHPort        int    `yaml:"ssh_port"`
	ContainerImage string `yaml:"container_image"`
	DataDir        string `yaml:"data_dir"`
	Variable       struct {
		Target string `yaml:"target"`
	} `yaml:"variable"`
}

type MasterConfig struct {
	Config struct {
		Listen  int    `yaml:"listen"`
		Prof    int    `yaml:"prof"`
		DataDir string `yaml:"data_dir"`
	} `yaml:"config"`
}

type MetaNodeConfig struct {
	Config struct {
		Listen  int    `yaml:"listen"`
		Prof    int    `yaml:"prof"`
		DataDir string `yaml:"data_dir"`
	} `yaml:"config"`
}

type DataNodeConfig struct {
	Config struct {
		Listen  int    `yaml:"listen"`
		Prof    int    `yaml:"prof"`
		DataDir string `yaml:"data_dir"`
	} `yaml:"config"`
}

type DeployHostsListConfig struct {
	Master struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"master"`
	MetaNode struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"metanode"`
	DataNode []struct {
		Hosts string   `yaml:"hosts"`
		Disk  []string `yaml:"disk"`
	} `yaml:"datanode"`
}

func CubeFSConfig() {
	// 读取配置文件
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("无法读取配置文件:", err)
		return
	}

	// 解析配置文件
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("无法解析配置文件:", err)
		return
	}

	// 打印配置信息
	fmt.Println("Global SSH Port:", config.Global.SSHPort)
	fmt.Println("Global Container Image:", config.Global.ContainerImage)
	fmt.Println("Global Data Directory:", config.Global.DataDir)
	fmt.Println("Global Variable Target:", config.Global.Variable.Target)

	fmt.Println("Master Listen Port:", config.Master.Config.Listen)
	fmt.Println("Master Prof Port:", config.Master.Config.Prof)
	fmt.Println("Master Data Directory:", config.Master.Config.DataDir)

	fmt.Println("MetaNode Listen Port:", config.MetaNode.Config.Listen)
	fmt.Println("MetaNode Prof Port:", config.MetaNode.Config.Prof)
	fmt.Println("MetaNode Data Directory:", config.MetaNode.Config.DataDir)

	fmt.Println("DataNode Listen Port:", config.DataNode.Config.Listen)
	fmt.Println("DataNode Prof Port:", config.DataNode.Config.Prof)
	fmt.Println("DataNode Data Directory:", config.DataNode.Config.DataDir)

	fmt.Println("Deploy Hosts List:")
	fmt.Println("Master Hosts:", config.DeployHostsList.Master.Hosts)
	fmt.Println("MetaNode Hosts:", config.DeployHostsList.MetaNode.Hosts)
	fmt.Println("DataNode Hosts and Disks:")
	for _, datanode := range config.DeployHostsList.DataNode {
		fmt.Println("  Hosts:", datanode.Hosts)
		fmt.Println("  Disks:", datanode.Disk)
	}
}
