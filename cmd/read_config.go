package cmd

import (
	"encoding/json"
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
		Listen  string `yaml:"listen"`
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

type Master struct {
	ClusterName         string `json:"clusterName"`
	ID                  string `json:"id"`
	Role                string `json:"role"`
	IP                  string `json:"ip"`
	Listen              string `json:"listen"`
	Prof                string `json:"prof"`
	Peers               string `json:"peers"`
	RetainLogs          string `json:"retainLogs"`
	ConsulAddr          string `json:"consulAddr"`
	ExporterPort        int    `json:"exporterPort"`
	LogLevel            string `json:"logLevel"`
	LogDir              string `json:"logDir"`
	WALDir              string `json:"walDir"`
	StoreDir            string `json:"storeDir"`
	MetaNodeReservedMem string `json:"metaNodeReservedMem"`
	EBSAddr             string `json:"ebsAddr"`
	EBSServicePath      string `json:"ebsServicePath"`
}

type MetaNode struct {
	Role              string   `json:"role"`
	Listen            string   `json:"listen"`
	Prof              string   `json:"prof"`
	RaftHeartbeatPort string   `json:"raftHeartbeatPort"`
	RaftReplicaPort   string   `json:"raftReplicaPort"`
	ConsulAddr        string   `json:"consulAddr"`
	ExporterPort      int      `json:"exporterPort"`
	LogLevel          string   `json:"logLevel"`
	LogDir            string   `json:"logDir"`
	WarnLogDir        string   `json:"warnLogDir"`
	TotalMem          string   `json:"totalMem"`
	MetadataDir       string   `json:"metadataDir"`
	RaftDir           string   `json:"raftDir"`
	MasterAddr        []string `json:"masterAddr"`
}

type DataNode struct {
	Role               string   `json:"role"`
	Listen             string   `json:"listen"`
	Prof               string   `json:"prof"`
	RaftHeartbeat      string   `json:"raftHeartbeat"`
	RaftReplica        string   `json:"raftReplica"`
	RaftDir            string   `json:"raftDir"`
	ConsulAddr         string   `json:"consulAddr"`
	ExporterPort       int      `json:"exporterPort"`
	Cell               string   `json:"cell"`
	LogDir             string   `json:"logDir"`
	LogLevel           string   `json:"logLevel"`
	Disks              []string `json:"disks"`
	DiskIopsReadLimit  string   `json:"diskIopsReadLimit"`
	DiskIopsWriteLimit string   `json:"diskIopsWriteLimit"`
	DiskFlowReadLimit  string   `json:"diskFlowReadLimit"`
	DiskFlowWriteLimit string   `json:"diskFlowWriteLimit"`
	MasterAddr         []string `json:"masterAddr"`
	EnableSmuxConnPool bool     `json:"enableSmuxConnPool"`
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
	// 将Master配置写入master.json文件

	master := Master{
		ClusterName:         "chubaofs01",
		ID:                  "1",
		Role:                "master",
		IP:                  "192.168.0.11",
		Listen:              "",
		Prof:                "17020",
		Peers:               "1:192.168.0.11:17010,2:192.168.0.12:17010,3:192.168.0.13:17010",
		RetainLogs:          "20000",
		ConsulAddr:          "http://192.168.0.101:8500",
		ExporterPort:        9500,
		LogLevel:            "debug",
		LogDir:              "/cfs/log",
		WALDir:              "/cfs/data/wal",
		StoreDir:            "/cfs/data/store",
		MetaNodeReservedMem: "67108864",
		EBSAddr:             "10.177.40.215:8500",
		EBSServicePath:      "access",
	}

	masterData, err := json.MarshalIndent(master, "", "  ")
	if err != nil {
		fmt.Println("无法编码Master配置 :", err)
		return
	}

	err = ioutil.WriteFile("master.json", masterData, 0644)
	if err != nil {
		fmt.Println("无法写入master.json文件:", err)
		return
	}

	// 将DataNode配置写入DataNode.json文件
	datanode := DataNode{
		Role:               "datanode",
		Listen:             "17310",
		Prof:               "17320",
		RaftHeartbeat:      "17330",
		RaftReplica:        "17340",
		RaftDir:            "/cfs/log",
		ConsulAddr:         "http://192.168.0.101:8500",
		ExporterPort:       9500,
		Cell:               "cell-01",
		LogDir:             "/cfs/log",
		LogLevel:           "debug",
		Disks:              []string{"/cfs/disk:10737418240"},
		DiskIopsReadLimit:  "20000",
		DiskIopsWriteLimit: "5000",
		DiskFlowReadLimit:  "1024000000",
		DiskFlowWriteLimit: "524288000",
		MasterAddr: []string{
			"192.168.0.11:17010",
			"192.168.0.12:17010",
			"192.168.0.13:17010",
		},
		EnableSmuxConnPool: true,
	}

	dataNodeData, err := json.MarshalIndent(datanode, "", "  ")
	if err != nil {
		fmt.Println("无法编码DataNode配置:", err)
		return
	}
	err = ioutil.WriteFile("dataNode.json", dataNodeData, 0644)
	if err != nil {
		fmt.Println("无法写入DataNode.json文件:", err)
		return
	}

	// 将MetaNode配置写入metanode.json文件
	metanode := MetaNode{
		Role:              "metanode",
		Listen:            "17210",
		Prof:              "17220",
		RaftHeartbeatPort: "17230",
		RaftReplicaPort:   "17240",
		ConsulAddr:        "http://192.168.0.101:8500",
		ExporterPort:      9500,
		LogLevel:          "debug",
		LogDir:            "/cfs/log",
		WarnLogDir:        "/cfs/log",
		TotalMem:          "536870912",
		MetadataDir:       "/cfs/data/meta",
		RaftDir:           "/cfs/data/raft",
		MasterAddr: []string{
			"192.168.0.11:17010",
			"192.168.0.12:17010",
			"192.168.0.13:17010",
		},
	}

	metaNodeData, err := json.MarshalIndent(metanode, "", "  ")
	if err != nil {
		fmt.Println("无法编码MetaNode配置:", err)
		return
	}
	err = ioutil.WriteFile("metanode.json", metaNodeData, 0644)
	if err != nil {
		fmt.Println("无法写入metanode.json文件:", err)
		return
	}
}
