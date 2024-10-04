package resource_type

import (
	"context"
	prick "prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
)

var _ interfaces.Prickable = &PrickableSynapseWorkspace{}

type PrickableSynapseWorkspace armsynapse.Workspace

type ListSynapseWorkspacesOptions struct {
	ResourceGroup string
}

func ListSynapseWorkspaces(api *prick.Api, options *ListSynapseWorkspacesOptions) ([]*PrickableSynapseWorkspace, error) {
	client := api.SynapseWorkspace
	var workspaces []*PrickableSynapseWorkspace
	switch options.ResourceGroup {
	case "":
		pager := client.NewListPager(&armsynapse.WorkspacesClientListOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, server := range res.Value {
				newServer := PrickableSynapseWorkspace(*server)
				workspaces = append(workspaces, &newServer)
			}
		}
	default:
		pager := client.NewListByResourceGroupPager(options.ResourceGroup, &armsynapse.WorkspacesClientListByResourceGroupOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, workspace := range res.Value {
				newSw := PrickableSynapseWorkspace(*workspace)
				workspaces = append(workspaces, &newSw)
			}
		}
	}

	return workspaces, nil
}

func (s *PrickableSynapseWorkspace) GetName() string {
	return *s.Name
}

func (s *PrickableSynapseWorkspace) GetLocation() string {
	return *s.Location
}

func (s *PrickableSynapseWorkspace) GetType() common.ResourceType {
	return common.ResourceTypeSynapseWorkspace
}

func (s *PrickableSynapseWorkspace) Poke(api *prick.Api) error {
	client := api.SynapseFirewall
	rg, err := common.ExtractResourceGroup(s.ID)
	if err != nil {
		return err
	}

	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}

	req, err := client.BeginCreateOrUpdate(
		context.Background(),
		*rg,
		*s.Name,
		ip,
		armsynapse.IPFirewallRuleInfo{
			Properties: &armsynapse.IPFirewallRuleProperties{
				StartIPAddress: &ip,
				EndIPAddress:   &ip,
			},
		},
		&armsynapse.IPFirewallRulesClientBeginCreateOrUpdateOptions{},
	)
	if err != nil {
		return err
	}

	_, err = req.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{Frequency: time.Duration(5 * float64(time.Second))})
	if err != nil {
		return err
	}

	return nil
}

func (s *PrickableSynapseWorkspace) Patch(api *prick.Api) error {
	client := api.SynapseFirewall
	rg, err := common.ExtractResourceGroup(s.ID)
	if err != nil {
		return err
	}

	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}

	req, err := client.BeginDelete(context.Background(), *rg, *s.Name, ip, &armsynapse.IPFirewallRulesClientBeginDeleteOptions{})
	if err != nil {
		return err
	}

	_, err = req.PollUntilDone(context.Background(), &runtime.PollUntilDoneOptions{Frequency: time.Duration(5 * float64(time.Second))})
	if err != nil {
		return err
	}

	return nil
}

func (s *PrickableSynapseWorkspace) ListPokes(api *prick.Api) ([]*common.Poke, error) {
	rules, err := s.ListIPRules(api)
	if err != nil {
		return nil, err
	}

	pokes := []*common.Poke{}
	for _, rule := range rules {
		pokes = append(pokes, &common.Poke{Name: *rule.Name, StartIpAddress: *rule.Properties.StartIPAddress, EndIpAddress: *rule.Properties.EndIPAddress})
	}
	return pokes, nil
}

func (s *PrickableSynapseWorkspace) ListIPRules(api *prick.Api) ([]*armsynapse.IPFirewallRuleInfo, error) {
	client := api.SynapseFirewall
	rg, err := common.ExtractResourceGroup(s.ID)
	if err != nil {
		return nil, err
	}

	pager := client.NewListByWorkspacePager(*rg, *s.Name, &armsynapse.IPFirewallRulesClientListByWorkspaceOptions{})
	var rules []*armsynapse.IPFirewallRuleInfo
	for pager.More() {
		res, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}

		rules = append(rules, res.Value...)
	}

	return rules, nil
}
