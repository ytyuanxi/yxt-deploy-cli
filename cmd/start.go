package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ip string
var allStart bool
var datanodeDisk string

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "start service",
	Long:  `This command will start services.`,
	Run: func(cmd *cobra.Command, args []string) {
		if allStart {
			fmt.Println("start all master services from config.yaml")
		} else {
			fmt.Println(cmd.UsageString())
		}

	},
}

var startMasterCommand = &cobra.Command{
	Use:   "master",
	Short: "start master",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("start master in ", ip)
		} else {
			fmt.Println("start all master services from config.yaml")
		}
	},
}

var startMetanodeCommand = &cobra.Command{
	Use:   "metanode",
	Short: "start metanode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("start metanode in ", ip)
		} else {
			fmt.Println("start all metanode services from config.yaml")
		}

	},
}

var startDatanodeCommand = &cobra.Command{
	Use:   "datanode",
	Short: "start datanode",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("start datanode in ", ip)
			if !cmd.Flags().Changed("disk") {
				fmt.Println("must have disk argument")
				os.Exit(1)
			}
			fmt.Println("disk:", datanodeDisk)
		} else {
			fmt.Println("start all datanode services from config.yaml")
		}
	},
}

func init() {
	StartCmd.AddCommand(startMasterCommand)
	StartCmd.AddCommand(startMetanodeCommand)
	StartCmd.AddCommand(startDatanodeCommand)
	StartCmd.Flags().BoolVarP(&allStart, "all", "a", false, "start all services")
	StartCmd.PersistentFlags().StringVarP(&ip, "ip", "", "", "specify an IP address to start services")
	startDatanodeCommand.Flags().StringVarP(&datanodeDisk, "disk", "d", "", "specify the disk where datanode mount")
}
