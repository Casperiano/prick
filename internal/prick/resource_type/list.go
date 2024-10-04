package resource_type

import (
	"fmt"
	prick "prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"
)

type ListPrickablesOptions struct {
	ResourceGroup string
	ResourceType  common.ResourceType
}

func ListPrickables(api *prick.Api, opts *ListPrickablesOptions) ([]interfaces.Prickable, error) {
	prickables := []interfaces.Prickable{}

	switch opts.ResourceType {
	case common.ResourceTypeStorageAccount:
		sa, err := ListStorageAccounts(api, &ListStorageAccountsOptions{ResourceGroup: opts.ResourceGroup})
		if err != nil {
			return nil, err
		}
		for _, sa := range sa {
			prickables = append(prickables, sa)

		}
	case common.ResourceTypeKeyVault:
		kv, err := ListKeyVaults(api, &ListKeyVaultsOptions{ResourceGroup: opts.ResourceGroup})
		if err != nil {
			return nil, err
		}
		for _, kv := range kv {
			prickables = append(prickables, kv)
		}
	case common.ResourceTypeSQLServer:
		ss, err := ListSQLServers(api, &ListSQLServersOptions{ResourceGroup: opts.ResourceGroup})
		if err != nil {
			return nil, err
		}
		for _, ss := range ss {
			prickables = append(prickables, ss)
		}
	case common.ResourceTypeSynapseWorkspace:
		sw, err := ListSynapseWorkspaces(api, &ListSynapseWorkspacesOptions{ResourceGroup: opts.ResourceGroup})
		if err != nil {
			return nil, err
		}
		for _, sw := range sw {
			prickables = append(prickables, sw)
		}
	default:
		return nil, fmt.Errorf("unsupported resource type: %v", opts.ResourceType)
	}
	return prickables, nil

}
