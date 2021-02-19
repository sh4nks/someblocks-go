package cmd

import (
	"os"
	"someblocks/internal/config"
	"someblocks/internal/server"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func runServer(cfg *config.Config) {
	srv := server.New(cfg)
	if srv == nil {
		os.Exit(1)
	}
	srv.Start()
}

func serverCmd(cfg *config.Config) *cobra.Command {
	var srvCmd = &cobra.Command{
		Use:   "server",
		Short: "Runs the webserver",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(cfg)
		},
	}

	srvCmd.PersistentFlags().IntP("port", "", 8080, "The port to bind to")
	srvCmd.PersistentFlags().StringP("host", "", "127.0.0.1", "The address to bind to")
	viper.BindPFlag("web.port", srvCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("web.host", srvCmd.PersistentFlags().Lookup("host"))
	return srvCmd
}
