package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Run unit test",
	Long:  `This command will run the unit tests for CubeFS.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running unit testsfor CubeFS…")
		// add your unit test logic here
	},
}
