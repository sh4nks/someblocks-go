package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"someblocks/config"
	"someblocks/utils"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var cfgFile string
var userLicense string

var rootCmd = &cobra.Command{
	Use:   "someblocks",
	Short: "The cli for managing my personal CMS written in golang",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("someblocks CMS v0.1-dev")
	},
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(func() {
		config.Load(cfgFile)
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: config.yaml)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd(config.Cfg))
	rootCmd.AddCommand(dbCmd(config.Cfg))
	rootCmd.AddCommand(userCmd(config.Cfg))
	rootCmd.AddCommand(configCmd(config.Cfg))

	noColor := false

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: noColor}
	output.FormatLevel = func(i interface{}) string {
		return utils.ColorizedFormatLevel(i, noColor)
	}
	log.Logger = log.Output(output)
}
