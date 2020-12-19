package cmd

import (
	"fmt"
	"someblocks/internal/config"

	"github.com/spf13/cobra"
)

var dumpCfg bool

func configCmd(cfg *config.Config) *cobra.Command {
	var cfgCmd = &cobra.Command{
		Use:   "config",
		Short: "Runs config related actions",
		Run: func(cmd *cobra.Command, args []string) {
			if dumpCfg {
				fmt.Printf(config.ToYAML())
			}
		},
	}
	cfgCmd.Flags().BoolVarP(&dumpCfg, "dump", "d", false, "dump the config")
	return cfgCmd
}
