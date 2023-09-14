package cmd

import (
	"io/ioutil"
	"log"
	"strconv"

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

	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Println("Unable to read configuration file:", err)
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("Unable to parse configuration file:", err)
		return nil, err
	}
	return config, nil

}

func readConfigTest() (*Config, error) {
	data, err := ioutil.ReadFile("config_test.yaml")
	if err != nil {
		log.Println("Unable to read configuration file:", err)
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println("Unable to parse configuration file:", err)
		return nil, err
	}
	return config, nil
}

type ConfigTest struct {
	Global   GlobalConfigTest   `yaml:"global"`
	Master   MasterConfigTest   `yaml:"master"`
	Metanode MetanodeConfigTest `yaml:"metanode"`
	Datanode DatanodeConfigTest `yaml:"datanode"`
}

type GlobalConfigTest struct {
	SSHPort        int    `yaml:"ssh_port"`
	ContainerImage string `yaml:"container_image"`
	DataDir        string `yaml:"data_dir"`
	BinDir         string `yaml:"bin_dir"`
	IP             string `yaml:"ip"`
	Variable       struct {
		Target string `yaml:"target"`
	} `yaml:"variable"`
}

type MasterConfigTest struct {
	Config []struct {
		Listen int `yaml:"listen"`
		Prof   int `yaml:"prof"`
	} `yaml:"config"`
}

type MetanodeConfigTest struct {
	Config []struct {
		Listen int `yaml:"listen"`
		Prof   int `yaml:"prof"`
	} `yaml:"config"`
}

type DatanodeConfigTest struct {
	Config []struct {
		Listen int `yaml:"listen"`
		Prof   int `yaml:"prof"`
		Disk   []struct {
			Path string `yaml:"path"`
			Size int64  `yaml:"size"`
		} `yaml:"disk"`
	} `yaml:"config"`
}

func convertToJosn() error {
	config, err := readConfig()
	if err != nil {
		return err
	}

	for id, node := range config.DeployHostsList.Master.Hosts {
		peers := getMasterPeers(config)
		err := writeMaster(ClusterName, strconv.Itoa(id+1), node, config.Master.Config.Listen, config.Master.Config.Prof, peers)
		if err != nil {
			return err
		}
	}

	masterAddr, err := getMasterAddrAndPort()
	if err != nil {
		return err
	}

	err = writeMetaNode(config.MetaNode.Config.Listen, config.MetaNode.Config.Prof, masterAddr)
	if err != nil {
		return err
	}

	disksInfo := []string{}
	for _, node := range config.DeployHostsList.DataNode {
		diskMap := ""
		for _, info := range node.Disk {
			diskMap += " -v " + info.Path + ":/cfs" + info.Path
			disksInfo = append(disksInfo, "/cfs"+info.Path+":"+info.Size)
		}
	}
	err = writeDataNode(config.DataNode.Config.Listen, config.DataNode.Config.Prof, masterAddr, disksInfo)
	if err != nil {
		return err
	}

	return nil
}
