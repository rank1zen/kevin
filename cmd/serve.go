package cmd

import (
	"context"
	"os"

	"github.com/rank1zen/kevin/internal/runtime"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start lol-service",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		cfg, err := runtime.NewConfig()
		if err != nil {
			panic(err)
		}

		rt, err := runtime.NewRuntime(ctx, cfg)
		if err != nil {
			panic(err)
		}

		os.Exit(rt.Run(ctx))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
