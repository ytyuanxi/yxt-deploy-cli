package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

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

func writeDataNode(listen, prof string, masterAddrs, disks []string) error {
	// 将DataNode配置写入DataNode.json文件
	datanode := DataNode{
		Role:               "datanode",
		Listen:             listen,
		Prof:               prof,
		RaftHeartbeat:      "17330",
		RaftReplica:        "17340",
		RaftDir:            "/cfs/log",
		ConsulAddr:         "http://192.168.0.101:8500",
		ExporterPort:       9500,
		Cell:               "cell-01",
		LogDir:             "/cfs/log",
		LogLevel:           "debug",
		Disks:              disks,
		DiskIopsReadLimit:  "20000",
		DiskIopsWriteLimit: "5000",
		DiskFlowReadLimit:  "1024000000",
		DiskFlowWriteLimit: "524288000",
		MasterAddr:         masterAddrs,
		EnableSmuxConnPool: true,
	}

	dataNodeData, err := json.MarshalIndent(datanode, "", "  ")
	if err != nil {
		log.Println("Unable to encode DataNode configuration:", err)
		return err
	}

	err = ioutil.WriteFile(ConfDir+"/datanode.json", dataNodeData, 0644)
	if err != nil {
		log.Println("Unable to write to DataNode.json file:", err)
		return err
	}
	return nil
}

func startAllDataNode() error {
	config, err := readConfig()
	if err != nil {
		log.Println(err)
	}
	masterAddr, err := getMasterAddrAndPort()
	if err != nil {
		return err
	}
	for id, node := range config.DeployHostsList.DataNode {

		disksInfo := []string{}
		diskMap := ""
		for _, info := range node.Disk {
			diskMap += " -v " + info.Path + ":/cfs" + info.Path
			disksInfo = append(disksInfo, "/cfs"+info.Path+":"+info.Size)
		}

		err := writeDataNode(config.DataNode.Config.Listen, config.DataNode.Config.Prof, masterAddr, disksInfo)
		if err != nil {
			return err
		}
		confFilePath := ConfDir + "/" + "datanode.json"
		err = transferConfigFileToRemote(confFilePath, config.Global.DataDir+"/"+ConfDir, RemoteUser, node.Hosts)
		if err != nil {
			return err
		}

		err = checkAndDeleteContainerOnNode(RemoteUser, node.Hosts, "datanode"+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		status, err := startDatanodeContainerOnNode(RemoteUser, node.Hosts, "datanode"+strconv.Itoa(id+1), config.Global.DataDir, diskMap)
		if err != nil {
			return err
		}
		log.Println(status)
	}
	log.Println("start all datanode services")
	return nil
}
func startDatanodeInSpecificNode(node string) error {

	config, err := readConfig()
	if err != nil {
		return err
	}
	for id, n := range config.DeployHostsList.DataNode {
		if n.Hosts == node {
			confFilePath := ConfDir + "/" + "datanode.json"
			err = transferDirectoryToRemote(confFilePath, config.Global.DataDir, RemoteUser, node)
			if err != nil {
				return err
			}

			err = checkAndDeleteContainerOnNode(RemoteUser, node, DataNodeName+strconv.Itoa(id+1))
			if err != nil {
				return err
			}
			diskMap := ""
			for _, info := range n.Disk {
				diskMap += " -v " + info.Path + ":/cfs" + info.Path

			}
			status, err := startDatanodeContainerOnNode(RemoteUser, node, DataNodeName+strconv.Itoa(id+1), config.Global.DataDir, diskMap)
			if err != nil {
				return err
			}

			log.Println(status)
			break
		}
	}
	return nil
}

func stopDatanodeInSpecificNode(node string) error {
	config, err := readConfig()
	if err != nil {
		return err
	}
	for id, n := range config.DeployHostsList.DataNode {
		if n.Hosts == node {
			status, err := stopContainerOnNode(RemoteUser, node, DataNodeName+strconv.Itoa(id+1))
			if err != nil {
				return err
			}
			log.Println(status)
			status, err = rmContainerOnNode(RemoteUser, node, DataNodeName+strconv.Itoa(id+1))
			if err != nil {
				return err
			}
			log.Println(status)
		}
	}
	return nil
}

func stopAllDataNode() error {
	config, err := readConfig()
	if err != nil {
		log.Println(err)
	}
	for id, node := range config.DeployHostsList.Master.Hosts {
		status, err := stopContainerOnNode(RemoteUser, node, DataNodeName+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		log.Println(status)
		status, err = rmContainerOnNode(RemoteUser, node, DataNodeName+strconv.Itoa(id+1))
		if err != nil {
			return err
		}
		log.Println(status)
	}
	return nil
}
