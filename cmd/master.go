package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

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

func getMasterPeers(config *Config) string {
	peers := ""
	for id, node := range config.DeployHostsList.Master.Hosts {
		if id != len(config.DeployHostsList.Master.Hosts)-1 {
			peers = peers + strconv.Itoa(id+1) + ":" + node + ":" + config.Master.Config.Listen + ","
		} else {
			peers = peers + strconv.Itoa(id+1) + ":" + node + ":" + config.Master.Config.Listen
		}

	}
	return peers
}

func writeMaster(clusterName, id, ip, listen, prof, peers string) error {
	master := Master{
		ClusterName:         clusterName,
		ID:                  id,
		Role:                "master",
		IP:                  ip,
		Listen:              listen,
		Prof:                prof,
		Peers:               peers,
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
		return fmt.Errorf("cannot be resolved to master.json %v", err)
	}
	fileName := ConfDir + "/master" + id + ".json"
	err = ioutil.WriteFile(fileName, masterData, 0644)
	if err != nil {

		return fmt.Errorf("unable to write %s  %v", fileName, err)
	}
	return nil
}

func startAllMaster() error {
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

		confFilePath := ConfDir + "/" + "master" + strconv.Itoa(id+1) + ".json"
		err = transferDirectoryToRemote(confFilePath, config.Global.DataDir, RemoteUser, node)
		if err != nil {
			return err
		}
		err = checkAndDeleteContainerOnNode(RemoteUser, node, MasterName+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		status, err := startMasterContainerOnNode(RemoteUser, node, MasterName+strconv.Itoa(id+1), config.Global.DataDir)
		if err != nil {
			return err
		}
		log.Println(status)
	}
	log.Println("start all master services")
	return nil
}

func stopAllMaster() error {
	config, err := readConfig()
	if err != nil {
		log.Println(err)
	}
	for id, node := range config.DeployHostsList.Master.Hosts {
		status, err := stopContainerOnNode(RemoteUser, node, MasterName+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		log.Println(status)
		status, err = rmContainerOnNode(RemoteUser, node, MasterName+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		log.Println(status)
	}
	return nil
}
