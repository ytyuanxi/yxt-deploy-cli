package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "cluster manager",
	Long:  `This command will manager the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		// add your cleanup logic here
	},
}

var subCommand = &cobra.Command{
	Use:   "init",
	Short: "init the cluster from config.yaml",
	Long:  "init the cluster from config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init-------------")
		ReadConfig()
	},
}

func init() {
	ClusterCmd.AddCommand(subCommand)
}
