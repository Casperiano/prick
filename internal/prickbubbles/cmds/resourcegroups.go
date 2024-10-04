package cmds

import (
	"prick/internal/prick"
	"prick/internal/prick/resource_group"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	tea "github.com/charmbracelet/bubbletea"
)

type ResourceGroupsFetchedMsg struct {
	ResourceGroups []*armresources.ResourceGroup
}

func FetchResourceGroups(api *prick.Api) tea.Cmd {
	return func() tea.Msg {
		resourceGroups, err := resource_group.ListResourceGroups(api)
		if err != nil {
			return err
		}

		return ResourceGroupsFetchedMsg{ResourceGroups: resourceGroups}
	}
}
