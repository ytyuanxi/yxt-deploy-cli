package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var allStop bool

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "start service",
	Long:  `This command will start service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if allStop {
			fmt.Println("stop all services......")
		} else {
			fmt.Println(cmd.UsageString())
		}
	},
}

var stopMasterCommand = &cobra.Command{
	Use:   "master",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("ip:", ip)
		} else {
			err := stopAllMaster()
			if err != nil {
				log.Println(err)
			}
			fmt.Println("stop all master services from config.yaml")
		}
	},
}

var stopMetanodeCommand = &cobra.Command{
	Use:   "metanode",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("stop metanode in ", ip)
		} else {
			//去目标节点stop该容器
			fmt.Println("stop all metanode services from config.yaml")
		}
	},
}

var stopDatanodeCommand = &cobra.Command{
	Use:   "datanode",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("stop datanode in ", ip)
			if !cmd.Flags().Changed("disk") {
				fmt.Println("must have disk argument")
				os.Exit(1)
			}
			fmt.Println("disk:", datanodeDisk)
		} else {
			fmt.Println("stop all datanode services from config.yaml")
		}
	},
}

func init() {
	StopCmd.AddCommand(stopMasterCommand)
	StopCmd.AddCommand(stopMetanodeCommand)
	StopCmd.AddCommand(stopDatanodeCommand)
	StopCmd.Flags().BoolVarP(&allStop, "all", "a", false, "stop all services")
	StopCmd.PersistentFlags().StringVarP(&ip, "ip", "", "", "specify an IP address to start services")
	stopDatanodeCommand.Flags().StringVarP(&datanodeDisk, "disk", "d", "", "specify the disk where datanode mount")

}
