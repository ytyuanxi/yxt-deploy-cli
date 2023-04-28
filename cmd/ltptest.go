package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var LtptestCmd = &cobra.Command{
	Use:   "ltptest",
	Short: "Run LTP test",
	Long:  `This command will run the LTP test.`,
	Run: func(cmd *cobra.Command, args []string) {
		if DiskPath == "" {
			cmd.Usage()
			os.Exit(1)
		}
		fmt.Printf("Running ltp test with disk path: %s\n", DiskPath)
		// add your LTP test logic here
	},
}
