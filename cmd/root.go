package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	initConfig()
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config/config.yaml)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(dbCmd)
}

func initConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(filepath.Join(".", "config")) // look for config in the src config directory
	viper.AddConfigPath("/etc/someblocks") // look for config in the src config directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
