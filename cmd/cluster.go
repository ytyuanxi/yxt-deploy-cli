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
		fmt.Println(cmd.UsageString())
	},
}

var initCommand = &cobra.Command{
	Use:   "init",
	Short: "init the cluster from config.yaml",
	Long:  "init the cluster from config.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		ReadConfig()
	},
}

var infoCommand = &cobra.Command{
	Use:   "info",
	Short: "Display cluster information",
	Long:  "Display cluster information",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var clearCommand = &cobra.Command{
	Use:   "clear",
	Short: "Clear cluster files and information",
	Long:  "Clear cluster files and information",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	ClusterCmd.AddCommand(initCommand)
	ClusterCmd.AddCommand(infoCommand)
	ClusterCmd.AddCommand(clearCommand)
}
