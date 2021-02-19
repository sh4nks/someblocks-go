package cmd

import (
	"someblocks/config"
	"someblocks/database"

	"github.com/spf13/cobra"
)

func dbCmd(cfg *config.Config) *cobra.Command {
	var dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Runs the database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			database.SetupAndMigrate(cfg)
		},
	}
	return dbCmd
}
