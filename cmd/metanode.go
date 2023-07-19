package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

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

func writeMetaNode(listen, prof string, masterAddrs []string) error {
	metanode := MetaNode{
		Role:              "metanode",
		Listen:            listen,
		Prof:              prof,
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
		MasterAddr:        masterAddrs,
	}

	metaNodeData, err := json.MarshalIndent(metanode, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("conf/metanode.json", metaNodeData, 0644)
	if err != nil {
		return err
	}
	return nil

}

func startAllMetaNode() error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	for id, node := range config.DeployHostsList.MetaNode.Hosts {

		err := writeMetaNode(config.MetaNode.Config.Listen, config.Master.Config.Prof, config.DeployHostsList.Master.Hosts)
		if err != nil {
			return err
		}
		confFilePath := ConfDir + "/" + "metanode.json"
		err = transferFileToRemote(confFilePath, config.Global.DataDir, RemoteUser, node)
		if err != nil {
			return err
		}

		err = checkAndDeleteContainerOnNode(RemoteUser, node, "metanode"+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		status, err := startMetanodeContainerOnNode(RemoteUser, node, "metanode"+strconv.Itoa(id+1), config.Global.DataDir)
		if err != nil {
			return err
		}
		log.Println(status)
	}
	log.Println("start all metanode services")
	return nil
}

func stopAllMetaNode() error {
	config, err := readConfig()
	if err != nil {
		log.Println(err)
	}
	for id, node := range config.DeployHostsList.Master.Hosts {
		//读取config，转化为对应的master.json
		peers := getMasterPeers(config)
		log.Println(peers)
		status, err := stopContainerOnNode(RemoteUser, node, "master"+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		log.Println(status)
	}
	fmt.Println("stop all master services from config.yaml")
	return nil
}
