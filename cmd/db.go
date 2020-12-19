package cmd

import (
	"someblocks/internal/models"
	"someblocks/internal/config"

	"github.com/spf13/cobra"
)



func dbCmd(cfg *config.Config) *cobra.Command {
	var dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Runs the database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			models.SetupAndMigrate(cfg)
		},
	}
	return dbCmd
}
