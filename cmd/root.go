package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"someblocks/config"
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
	cobra.OnInitialize(func() { config.Load(cfgFile) })

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config/config.yaml)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(dbCmd)
}
