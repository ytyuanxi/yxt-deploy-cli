package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type MetaNode struct {
	Role              string   `json:"role"`
	Listen            string   `json:"listen"`
	Prof              string   `json:"prof"`
	RaftHeartbeatPort string   `json:"raftHeartbeatPort"`
	RaftReplicaPort   string   `json:"raftReplicaPort"`
	LocalIP           string   `json:"localIP"`
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

func readMetaNode(filename string) (*MetaNode, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	metaNode := &MetaNode{}
	err = json.Unmarshal(data, metaNode)
	if err != nil {
		return nil, err
	}
	return metaNode, nil
}

func writeMetaNode(listen, prof, id, localIP string, masterAddrs []string) error {
	metanode := MetaNode{
		Role:              "metanode",
		Listen:            listen,
		Prof:              prof,
		RaftHeartbeatPort: "17230",
		RaftReplicaPort:   "17240",
		LocalIP:           localIP,
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
	err = ioutil.WriteFile("conf/metanode"+id+".json", metaNodeData, 0644)
	if err != nil {
		return err
	}
	return nil

}

func stopMetanodeInSpecificNode(node string) error {
	//获取该ip对应的容器名

	files, err := ioutil.ReadDir(ConfDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "metanode") && !file.IsDir() {
			data, err := readMetaNode(ConfDir + "/" + file.Name())
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file.Name(), err)
				return nil
			}
			if data.LocalIP == node {
				status, err := stopContainerOnNode(RemoteUser, node, strings.Split(file.Name(), ".")[0])
				if err != nil {
					return err
				}
				log.Println(status)
				status, err = rmContainerOnNode(RemoteUser, node, strings.Split(file.Name(), ".")[0])
				if err != nil {
					return err
				}
				log.Println(status)
			}
		}
	}

	return nil

}

// func stopMetanodeInSpecificNode(node string) error {
// 	//获取该ip对应的容器名
// 	config, err := readConfig()
// 	if err != nil {
// 		return err
// 	}
// 	for id, n := range config.DeployHostsList.MetaNode.Hosts {
// 		if node == n {
// 			status, err := stopContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1))
// 			if err != nil {
// 				return err
// 			}
// 			log.Println(status)
// 			status, err = rmContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1))
// 			if err != nil {
// 				return err
// 			}
// 			log.Println(status)
// 		}

// 	}
// 	return nil

// }

func startMetanodeInSpecificNode(node string) error {
	//找到对应ip的配置文件
	config, err := readConfig()
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(ConfDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "metanode") && !file.IsDir() {
			data, err := readMetaNode(ConfDir + "/" + file.Name())
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file.Name(), err)
				return nil
			}
			if data.LocalIP == node {
				var dataDir string
				if config.Master.Config.DataDir == "" {
					dataDir = config.Global.DataDir
				} else {
					dataDir = config.Master.Config.DataDir

				}
				confFilePath := ConfDir + "/" + file.Name()
				err = transferConfigFileToRemote(confFilePath, dataDir+"/"+ConfDir, RemoteUser, node)
				if err != nil {
					return err
				}

				err = checkAndDeleteContainerOnNode(RemoteUser, node, strings.Split(file.Name(), ".")[0])
				if err != nil {
					return err
				}
				status, err := startMetanodeContainerOnNode(RemoteUser, node, strings.Split(file.Name(), ".")[0], dataDir)
				if err != nil {
					return err
				}
				log.Println(status)
				break
			}

		}
	}

	return nil
}

// func startMetanodeInSpecificNode(node string) error {
// 	//要对执行ip启动的容器进行编号
// 	config, err := readConfig()
// 	if err != nil {
// 		return err
// 	}

// 	var dataDir string
// 	if config.Master.Config.DataDir == "" {
// 		dataDir = config.Global.DataDir
// 	} else {
// 		dataDir = config.Master.Config.DataDir
// 	}
// 	for id, n := range config.DeployHostsList.MetaNode.Hosts {
// 		if n == node {
// 			confFilePath := ConfDir + "/" + "metanode.json"

// 			err = transferConfigFileToRemote(confFilePath, dataDir+"/"+ConfDir, RemoteUser, node)
// 			if err != nil {
// 				return err
// 			}

// 			err = checkAndDeleteContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1))
// 			if err != nil {
// 				return err
// 			}
// 			status, err := startMetanodeContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1), dataDir)
// 			if err != nil {
// 				return err
// 			}
// 			log.Println(status)
// 			break
// 		}
// 	}
// 	return nil
// }

func getMasterAddrAndPort() ([]string, error) {
	config, err := readConfig()
	if err != nil {
		return []string{}, err
	}
	masterAddr := make([]string, len(config.DeployHostsList.Master.Hosts))
	for id, node := range config.DeployHostsList.Master.Hosts {
		masterAddr[id] = node + ":" + config.Master.Config.Listen
	}
	return masterAddr, nil
}

func startAllMetaNode() error {

	config, err := readConfig()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(ConfDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "metanode") && !file.IsDir() {
			data, err := readMetaNode(ConfDir + "/" + file.Name())
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file.Name(), err)
				return nil
			}
			confFilePath := ConfDir + "/" + file.Name()
			var dataDir string
			if config.Master.Config.DataDir == "" {
				dataDir = config.Global.DataDir
			} else {
				dataDir = config.Master.Config.DataDir

			}
			err = transferConfigFileToRemote(confFilePath, dataDir+"/"+ConfDir, RemoteUser, data.LocalIP)
			if err != nil {
				return err
			}
			err = checkAndDeleteContainerOnNode(RemoteUser, data.LocalIP, strings.Split(file.Name(), ".")[0])
			if err != nil {
				return err
			}
			status, err := startMetanodeContainerOnNode(RemoteUser, data.LocalIP, strings.Split(file.Name(), ".")[0], dataDir)
			if err != nil {
				return err
			}
			log.Println(status)
		}
	}

	// confFilePath := ConfDir + "/" + "metanode.json"
	// var dataDir string
	// if config.Master.Config.DataDir == "" {
	// 	dataDir = config.Global.DataDir
	// } else {
	// 	dataDir = config.Master.Config.DataDir
	// }
	// for id, node := range config.DeployHostsList.MetaNode.Hosts {
	// 	err = transferConfigFileToRemote(confFilePath, dataDir+"/"+ConfDir, RemoteUser, node)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	err = checkAndDeleteContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1))
	// 	if err != nil {
	// 		return err
	// 	}
	// 	status, err := startMetanodeContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1), dataDir)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println(status)
	// }

	//Detect successful deployment
	log.Println("start all metanode services")
	return nil
}

// func startAllMetaNode() error {
// 	config, err := readConfig()
// 	if err != nil {
// 		return err
// 	}
// 	// masterAddr, err := getMasterAddrAndPort()
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	for id, node := range config.DeployHostsList.MetaNode.Hosts {

// 		// err := writeMetaNode(config.MetaNode.Config.Listen, config.MetaNode.Config.Prof, masterAddr)
// 		// if err != nil {
// 		// 	return err
// 		// }
// 		confFilePath := ConfDir + "/" + "metanode.json"
// 		err = transferDirectoryToRemote(confFilePath, config.Global.DataDir+"/"+ConfDir, RemoteUser, node)
// 		if err != nil {
// 			return err
// 		}

// 		err = checkAndDeleteContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1))
// 		if err != nil {
// 			return err
// 		}
// 		status, err := startMetanodeContainerOnNode(RemoteUser, node, MetaNodeName+strconv.Itoa(id+1), config.Global.DataDir)
// 		if err != nil {
// 			return err
// 		}
// 		log.Println(status)
// 	}
// 	log.Println("start all metanode services")
// 	return nil
// }

func stopAllMetaNode() error {

	files, err := ioutil.ReadDir(ConfDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "metanode") && !file.IsDir() {
			data, err := readMetaNode(ConfDir + "/" + file.Name())
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", file.Name(), err)
				return nil
			}
			status, err := stopContainerOnNode(RemoteUser, data.LocalIP, strings.Split(file.Name(), ".")[0])
			if err != nil {
				return err
			}
			log.Println(status)
			status, err = rmContainerOnNode(RemoteUser, data.LocalIP, strings.Split(file.Name(), ".")[0])
			if err != nil {
				return err
			}
			log.Println(status)
		}
	}
	return nil
}
