package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start CubeFS client Docker image",
	Long:  `This command will start the CubeFS client Docker image.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting CubeFS client...")
		// add your client start logic here
	},
}
