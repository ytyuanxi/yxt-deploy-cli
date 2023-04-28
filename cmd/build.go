package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build CubeFS server and client",
	Long:  `Builds CubeFS server and client.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := exec.Command("bash", "-c", "bash /path/to/cubefs/build.sh").Run(); err != nil {
			fmt.Println("Error building CubeFS server and client:", err)
			os.Exit(1)
		}
		fmt.Println("CubeFS server and client built successfully")
	},
}
