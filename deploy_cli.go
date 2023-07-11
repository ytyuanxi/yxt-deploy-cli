package main

import (
	"deploy_cli/cmd"
	"fmt"
)

func init() {
	cmd.RootCmd.Flags().BoolVarP(&cmd.Version, "version", "v", false, "show version information")
	cmd.RootCmd.AddCommand(cmd.ClusterCmd)
	cmd.RootCmd.AddCommand(cmd.StartCmd)
	cmd.RootCmd.AddCommand(cmd.StopCmd)
	cmd.RootCmd.AddCommand(cmd.RestartCmd)

}

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
}
