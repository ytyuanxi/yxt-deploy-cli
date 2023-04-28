package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var ScenariotestCmd = &cobra.Command{
	Use:   "scenariotest",
	Short: "Run scenario test",
	Long:  `This command will run the scenario test.`,
	Run: func(cmd *cobra.Command, args []string) {
		if DiskPath == "" {
			cmd.Usage()
			os.Exit(1)
		}
		fmt.Printf("Running scenario test with disk path: %s\n", DiskPath)
		// add your scenario test logic here
	},
}
