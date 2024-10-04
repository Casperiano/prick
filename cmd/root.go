package cmd

import (
	"context"
	"os"

	"prick/internal/prick/common"
	"prick/internal/prick/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "v0.0.1",
	Use:     "prick",
	Short:   "Prick is a CLI tool for developers to manage firewall rules on Azure resources.",
}

var (
	flagResourceGroup string
	flagResource      string
	flagResourceType  common.ResourceType
)

type contextKey string

func Execute() error {
	return rootCmd.ExecuteContext(context.Background())
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		config, err := config.LoadConfig()
		if err != nil {
			return err
		}

		ctx := context.WithValue(cmd.Context(), contextKey("config"), config)
		cmd.SetContext(ctx)

		return nil
	}
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			err := cmd.Help()
			if err != nil {
				return err
			}
			os.Exit(0)
		}

		return nil
	}
}
