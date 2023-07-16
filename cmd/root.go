package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version bool

const ClusterName = "cubeFS"
const RemoteUser = "root"

var RootCmd = &cobra.Command{
	Use:   "deploy-cli",
	Short: "CLI for managing CubeFS server and client using Docker",
	Long:  `cubefs is a CLI application for managing CubeFS, an open-source distributed file system, using Docker containers.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			fmt.Println("0.0.1")
		} else {
			fmt.Println(cmd.UsageString())
		}
	},
}
