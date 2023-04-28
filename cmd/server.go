package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start CubeFS servers Docker image",
	Long:  `This command will start the CubeFS servers Docker image.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting CubeFS servers...")
		// add your server start logic here
	},
}
