package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Cleanup old Docker image",
	Long:  `This command will cleanup old Docker image.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cleaning up old Docker image...")
		// add your cleanup logic here
	},
}
