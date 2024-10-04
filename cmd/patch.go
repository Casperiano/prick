package cmd

import (
	"errors"
	"fmt"
	prick "prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"
	"prick/internal/prick/resource_type"
	"sync"

	"github.com/spf13/cobra"
)

func init() {
	pokeCmd := &cobra.Command{
		Use:   "patch",
		Short: "",
		Long:  ``,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if (flagResource != "") && (flagResourceType == "") && (flagResourceGroup == "") {
				return errors.New("resource group or resource type is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			api, err := prick.NewApi()
			if err != nil {
				return err
			}
			fmt.Printf("🔎 Searching %s(s)...\n", flagResourceType)
			var resources []interfaces.Prickable

			switch flagResourceType {
			case common.ResourceTypeStorageAccount:
				sa, err := resource_type.ListStorageAccounts(api, &resource_type.ListStorageAccountsOptions{ResourceGroup: flagResourceGroup})
				resources = make([]interfaces.Prickable, len(sa))
				for i, s := range sa {
					resources[i] = s
				}
				if err != nil {
					return fmt.Errorf("error: %v", err)
				}
			case common.ResourceTypeKeyVault:
				kv, err := resource_type.ListKeyVaults(api, &resource_type.ListKeyVaultsOptions{ResourceGroup: flagResourceGroup})
				resources = make([]interfaces.Prickable, len(kv))
				for i, s := range kv {
					resources[i] = s
				}
				if err != nil {
					return fmt.Errorf("error: %v", err)
				}
			case common.ResourceTypeSQLServer:
				sql, err := resource_type.ListSQLServers(api, &resource_type.ListSQLServersOptions{ResourceGroup: flagResourceGroup})
				resources = make([]interfaces.Prickable, len(sql))
				for i, s := range sql {
					resources[i] = s
				}
				if err != nil {
					return fmt.Errorf("error: %v", err)
				}
			default:
				return fmt.Errorf("unsupported resource type: %s", flagResourceType)
			}

			wg := sync.WaitGroup{}

			for _, resource := range resources {
				if resource.GetName() == flagResource || flagResource == "" {
					wg.Add(1)
					ps := resource

					go func(wg *sync.WaitGroup) {
						defer wg.Done()
						err := ps.Patch(api)
						if err != nil {
							fmt.Printf("🧱 Error patching storage account %s: %v\n", ps.GetName(), err)
						} else {
							fmt.Printf("🩹 Successfully patched storage account %s\n", ps.GetName())
						}
					}(&wg)

				}
			}

			wg.Wait()

			return nil
		},
	}

	rootCmd.AddCommand(pokeCmd)
	pokeCmd.PersistentFlags().StringVar(&flagResourceGroup, "resource-group", "", "Resource group name")
	pokeCmd.PersistentFlags().Var(&flagResourceType, "resource-type", "Resource type")
	pokeCmd.PersistentFlags().StringVar(&flagResource, "resource", "", "Resource group name")
}
