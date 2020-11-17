package cmd

import (
	"someblocks/internal/app"

	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Runs the database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.App{}
		app.Migrate()
	},
}
