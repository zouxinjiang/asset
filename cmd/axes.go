package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zouxinjiang/axes/cmd/server"
)

var (
	axesCmd = &cobra.Command{
		Use:   "axes",
		Short: "axes",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
)

func Execute() error {
	return axesCmd.Execute()
}

func init() {
	axesCmd.AddCommand(server.ServerCmd)
}
