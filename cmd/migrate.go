package cmd

import (
	"context"
	"os"

	"github.com/rank1zen/kevin/internal/app"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate to current schema version",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		a := app.NewMigrator(ctx)
		os.Exit(a.Run(ctx))
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
