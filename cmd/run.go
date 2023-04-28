package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run servers, client and monitor",
	Long:  `This command will run the CubeFS servers, client and monitor.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running CubeFS servers, client and monitor...")
		// add your run logic here
	},
}
