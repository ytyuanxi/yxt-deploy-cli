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
	BinDir         string `yaml:"bin_dir"`
	Variable       struct {
		Target string `yaml:"target"`
	} `yaml:"variable"`
}

type MasterConfig struct {
	Config struct {
		Listen  string `yaml:"listen"`
		Prof    string `yaml:"prof"`
		DataDir string `yaml:"data_dir"`
	} `yaml:"config"`
}

type MetaNodeConfig struct {
	Config struct {
		Listen  string `yaml:"listen"`
		Prof    string `yaml:"prof"`
		DataDir string `yaml:"data_dir"`
	} `yaml:"config"`
}

type DataNodeConfig struct {
	Config struct {
		Listen  string `yaml:"listen"`
		Prof    string `yaml:"prof"`
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
		Hosts string     `yaml:"hosts"`
		Disk  []DiskInfo `yaml:"disk"`
	} `yaml:"datanode"`
}

type DiskInfo struct {
	Path string `yaml:"path"`
	Size string `yaml:"size"`
}

func readConfig() (*Config, error) {
	// 读取配置文件
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println("无法读取配置文件:", err)
		return nil, err
	}

	// 解析配置文件
	config := &Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("无法解析配置文件:", err)
		return nil, err
	}
	return config, nil

}

func tmp() {

	// 将MetaNode配置写入metanode.json文件

}
