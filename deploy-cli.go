package main

import (
	"deploy_cli/cmd"
	"fmt"
)

func init() {
	cmd.RootCmd.PersistentFlags().StringVarP(&cmd.DiskPath, "disk", "d", "", "set CubeFS DataNode local disk path")
	cmd.RootCmd.Flags().BoolVarP(&cmd.Version, "version", "v", false, "show version information")
	cmd.RootCmd.Flags().BoolVarP(&cmd.All, "all-servers", "a", false, "deploy all servers")

	cmd.RootCmd.AddCommand(cmd.BuildCmd)
	cmd.RootCmd.AddCommand(cmd.ServerCmd)
	cmd.RootCmd.AddCommand(cmd.ClientCmd)
	cmd.RootCmd.AddCommand(cmd.MonitorCmd)
	cmd.RootCmd.AddCommand(cmd.LtptestCmd)
	cmd.RootCmd.AddCommand(cmd.ScenariotestCmd)
	cmd.RootCmd.AddCommand(cmd.RunCmd)
	cmd.RootCmd.AddCommand(cmd.CleanCmd)
	cmd.RootCmd.AddCommand(cmd.TestCmd)
}

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
}
