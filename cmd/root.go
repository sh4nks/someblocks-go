package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"someblocks/internal/config"
	"someblocks/pkg/utils"

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
	rootCmd.AddCommand(configCmd(config.Cfg))

	noColor := false

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				l = utils.Colorize("TRACE", utils.ColorMagenta, noColor)
			case "debug":
				l = utils.Colorize("DEBUG", utils.ColorYellow, noColor)
			case "info":
				l = utils.Colorize("INFO", utils.ColorGreen, noColor)
			case "warn":
				l = utils.Colorize("WARN", utils.ColorRed, noColor)
			case "error":
				l = utils.Colorize(utils.Colorize("ERROR", utils.ColorRed, noColor), utils.ColorBold, noColor)
			case "fatal":
				l = utils.Colorize(utils.Colorize("FATAL", utils.ColorRed, noColor), utils.ColorBold, noColor)
			case "panic":
				l = utils.Colorize(utils.Colorize("PANIC", utils.ColorRed, noColor), utils.ColorBold, noColor)
			default:
				l = utils.Colorize("???", utils.ColorBold, noColor)
			}
		} else {
			if i == nil {
				l = utils.Colorize("???", utils.ColorBold, noColor)
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))[0:3]
			}
		}

		if noColor {
			return fmt.Sprintf(" %-6s", l)
		}
		// 14 - because terminal colors are taking up some bytes
		return fmt.Sprintf(" %-15s", l)
	}
	log.Logger = log.Output(output)
}
