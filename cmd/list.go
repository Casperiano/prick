package cmd

import (
	"fmt"
	"os"
	prick "prick/internal/prick"
	"prick/internal/prick/config"
	"prick/internal/prickbubbles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// if (flagResource != "") && (flagResourceType == "") && (flagResourceGroup == "") {
			// 	return errors.New("resource group or resource type is required")
			// }
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cmd.Context().Value(contextKey("config")).(*config.Config)
			api, err := prick.NewApi()
			if err != nil {
				return err
			}

			logEnabled := os.Getenv("PRICK_LOG")
			if logEnabled != "" {
				f, err := tea.LogToFile("debug.log", "debug")
				if err != nil {
					fmt.Println("fatal:", err)
					os.Exit(1)
				}
				defer f.Close()
			}

			model := prickbubbles.New(api, config)
			if _, err := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseAllMotion()).Run(); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(listCmd)
	listCmd.PersistentFlags().StringVar(&flagResourceGroup, "resource-group", "", "Resource group name")
	listCmd.PersistentFlags().Var(&flagResourceType, "resource-type", "Resource type")
}
