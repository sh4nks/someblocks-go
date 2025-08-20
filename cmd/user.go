package cmd

import (
	"someblocks/config"
	"someblocks/database"
	"someblocks/models"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	username string
	email    string
	role     string
	password string
	size     int
)

func userCmd(cfg *config.Config) *cobra.Command {
	var dbUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Modify a user",
	}
	dbUserCmd.AddCommand(userAddCmd(cfg))
	dbUserCmd.AddCommand(userListCmd(cfg))
	dbUserCmd.AddCommand(userUpdateCmd(cfg))
	return dbUserCmd
}

func userAddCmd(cfg *config.Config) *cobra.Command {
	var dbUserCmd = &cobra.Command{
		Use:   "add",
		Short: "Insert a new user",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := database.SetupDatabase(cfg)
			userService := models.UserService{DB: db}
			if err != nil {
				log.Error().Err(err).Msg("Couldn't setup database")
				return
			}
			userService.Insert(username, email, password)
		},
	}

	dbUserCmd.Flags().StringVar(&username, "username", "", "The username")
	dbUserCmd.Flags().StringVar(&password, "password", "", "The password")
	dbUserCmd.Flags().StringVar(&role, "role", "", "The role")
	dbUserCmd.Flags().StringVar(&email, "email", "", "The email")
	return dbUserCmd
}

func userListCmd(cfg *config.Config) *cobra.Command {

	var dbUserCmd = &cobra.Command{
		Use:   "search",
		Short: "List all users or search for a specific one",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := database.SetupDatabase(cfg)
			userService := models.UserService{DB: db}
			if err != nil {
				log.Error().Err(err).Msg("Couldn't setup database")
			} else {
				result := userService.Search(username, size)

				timeformat := "2006-01-02 15:04:05"

				log.Info().Msgf("found %d user(s)", len(result))
				for _, user := range result {
					log.Info().Msgf("id: %d, username: %s, email: %s, created: %s, updated: %s", user.ID, user.Username, user.Email, user.CreatedAt.Format(timeformat), user.UpdatedAt.Format(timeformat))
				}
			}
		},
	}

	dbUserCmd.Flags().StringVar(&username, "username", "", "The username")
	dbUserCmd.Flags().IntVar(&size, "size", -1, "The amount of users returned")
	return dbUserCmd
}

func userUpdateCmd(cfg *config.Config) *cobra.Command {
	var dbUserCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an existing user",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := database.SetupDatabase(cfg)
			userService := models.UserService{DB: db}
			if err != nil {
				log.Error().Err(err).Msg("Couldn't setup database")
			} else {
				userService.UpdatePassword(&models.User{Username: username}, password)
				log.Info().Msgf("password updated for user %s", username)
			}
		},
	}

	dbUserCmd.Flags().StringVar(&username, "username", "", "The username")
	dbUserCmd.Flags().StringVar(&password, "password", "", "The password")
	dbUserCmd.Flags().StringVar(&role, "role", "", "The role")
	dbUserCmd.Flags().StringVar(&email, "email", "", "The email")
	return dbUserCmd
}
