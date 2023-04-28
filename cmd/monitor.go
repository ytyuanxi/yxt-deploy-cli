package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Start monitor web UI",
	Long:  `This command will start the monitor web UI.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting monitor web UI...")
		// add your monitor start logic here
	},
}
