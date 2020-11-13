package cmd

import (
	"fmt"
	"log"
	"net/http"
	"someblocks/internal/app"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


func runServer(host string, port int) {
	app := app.App{}
	app.CreateApp()

	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Running on http://%s/ (Press CTRL+C to quit)", addr)

	http.ListenAndServe(addr, app.Routes)
}

func serverCmd() *cobra.Command {
	var srvCmd = &cobra.Command{
		Use:   "server",
		Short: "Runs the webserver",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(viper.GetString("web.host"), viper.GetInt("web.port"))
		},
	}

	srvCmd.PersistentFlags().IntP("port", "", 8080, "The port to bind to")
	srvCmd.PersistentFlags().StringP("host", "", "127.0.0.1", "The address to bind to")
	viper.BindPFlag("web.port", srvCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("web.host", srvCmd.PersistentFlags().Lookup("host"))
	return srvCmd
}
