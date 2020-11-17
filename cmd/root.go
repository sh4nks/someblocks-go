package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"someblocks/internal/config"
	"someblocks/internal/core"

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
	cobra.OnInitialize(func() { config.Load(cfgFile) })

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: config.yaml)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(serverCmd())
	rootCmd.AddCommand(dbCmd)

	// TODO: Bind to cobra and viper
	noColor := false

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "trace":
				l = core.Colorize("TRACE", core.ColorMagenta, noColor)
			case "debug":
				l = core.Colorize("DEBUG", core.ColorYellow, noColor)
			case "info":
				l = core.Colorize("INFO", core.ColorGreen, noColor)
			case "warn":
				l = core.Colorize("WARN", core.ColorRed, noColor)
			case "error":
				l = core.Colorize(core.Colorize("ERROR", core.ColorRed, noColor), core.ColorBold, noColor)
			case "fatal":
				l = core.Colorize(core.Colorize("FATAL", core.ColorRed, noColor), core.ColorBold, noColor)
			case "panic":
				l = core.Colorize(core.Colorize("PANIC", core.ColorRed, noColor), core.ColorBold, noColor)
			default:
				l = core.Colorize("???", core.ColorBold, noColor)
			}
		} else {
			if i == nil {
				l = core.Colorize("???", core.ColorBold, noColor)
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
