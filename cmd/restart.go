package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var allRestart bool

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "start service",
	Long:  `This command will start service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if allRestart {
			fmt.Println("restart all services......")
		} else {
			fmt.Println(cmd.UsageString())
		}
	},
}

var restartMasterCommand = &cobra.Command{
	Use:   "master",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("ip:", ip)
		} else {
			fmt.Println("restart all master services from config.yaml")
		}
	},
}

var restartMetanodeCommand = &cobra.Command{
	Use:   "metanode",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("ip:", ip)
		} else {
			fmt.Println("restart all metanode services from config.yaml")
		}

	},
}

var restartDatanodeCommand = &cobra.Command{
	Use:   "datanode",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("ip") {
			fmt.Println("ip:", ip)
			if !cmd.Flags().Changed("disk") {
				fmt.Println("must have disk argument")
				os.Exit(1)
			}
			fmt.Println("disk:", datanodeDisk)
		} else {
			fmt.Println("restart all datanode services from config.yaml")
		}
	},
}

func init() {
	RestartCmd.AddCommand(restartMasterCommand)
	RestartCmd.AddCommand(restartMetanodeCommand)
	RestartCmd.AddCommand(restartDatanodeCommand)
	RestartCmd.Flags().BoolVarP(&allRestart, "all", "a", false, "stop all services")
	RestartCmd.PersistentFlags().StringVarP(&ip, "ip", "", "", "specify an IP address to start services")
	restartDatanodeCommand.Flags().StringVarP(&datanodeDisk, "disk", "d", "", "specify the disk where datanode mount")

}
