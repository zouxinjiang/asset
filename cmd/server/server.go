package server

import (
	"github.com/spf13/cobra"
	"github.com/zouxinjiang/axes/internal/app/server"
	"github.com/zouxinjiang/axes/pkg/cobra_args_parser"
	"github.com/zouxinjiang/axes/pkg/errors"
)

var (
	ServerCmd = &cobra.Command{
		Use:   "server",
		Short: "axes server tool",
		Long:  "axes server tool",
		RunE:  runServer,
	}
)

func init() {
	cobra_args_parser.InitArgs(server.Args{}, ServerCmd)
}

func runServer(cmd *cobra.Command, args []string) error {
	appArgs := server.Args{}
	cobra_args_parser.ParseArgs(&appArgs, cmd)
	err := server.Run(appArgs)
	if err != nil {
		err = errors.Wrap(err, errors.CodeUnexpect)
		return err
	}
	return nil
}
