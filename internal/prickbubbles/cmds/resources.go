package cmds

import (
	"log"
	"prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"
	resource_type "prick/internal/prick/resource_type"

	tea "github.com/charmbracelet/bubbletea"
)

type ResourcesFetchedMsg struct {
	Prickables []interfaces.Prickable
}

func FetchResources(api *prick.Api, rg string) tea.Cmd {
	return func() tea.Msg {
		var prickables = []interfaces.Prickable{}
		for _, rt := range common.ResourceTypes() {
			ps, err := resource_type.ListPrickables(api, &resource_type.ListPrickablesOptions{
				ResourceGroup: rg,
				ResourceType:  rt,
			})
			if err != nil {
				log.Default().Printf("error listing prickables: %v", err)
			}
			prickables = append(prickables, ps...)
		}

		return ResourcesFetchedMsg{Prickables: prickables}
	}
}
