package resource_type

import (
	"context"
	prick "prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
)

var _ interfaces.Prickable = &PrickableSQLServer{}

type PrickableSQLServer armsql.Server

type ListSQLServersOptions struct {
	ResourceGroup string
}

func ListSQLServers(api *prick.Api, options *ListSQLServersOptions) ([]*PrickableSQLServer, error) {
	client := api.SqlServer
	var servers []*PrickableSQLServer
	switch options.ResourceGroup {
	case "":
		pager := client.NewListPager(&armsql.ServersClientListOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, server := range res.Value {
				newServer := PrickableSQLServer(*server)
				servers = append(servers, &newServer)
			}
		}
	default:
		pager := client.NewListByResourceGroupPager(options.ResourceGroup, &armsql.ServersClientListByResourceGroupOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, server := range res.Value {
				newKv := PrickableSQLServer(*server)
				servers = append(servers, &newKv)
			}
		}
	}

	return servers, nil
}

func (s *PrickableSQLServer) GetName() string {
	return *s.Name
}

func (s *PrickableSQLServer) GetLocation() string {
	return *s.Location
}

func (s *PrickableSQLServer) GetType() common.ResourceType {
	return common.ResourceTypeSQLServer
}

func (s *PrickableSQLServer) Poke(api *prick.Api) error {
	client := api.SqlFirewall
	rg, err := common.ExtractResourceGroup(s.ID)
	if err != nil {
		return err
	}

	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}

	_, err = client.CreateOrUpdate(
		context.Background(),
		*rg,
		*s.Name,
		"myIP",
		armsql.FirewallRule{Properties: &armsql.ServerFirewallRuleProperties{StartIPAddress: &ip, EndIPAddress: &ip}},
		&armsql.FirewallRulesClientCreateOrUpdateOptions{})
	return err
}

func (s *PrickableSQLServer) Patch(api *prick.Api) error {
	client := api.SqlFirewall
	rg, err := common.ExtractResourceGroup(s.ID)
	if err != nil {
		return err
	}

	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}
	ipRules, err := s.ListIPRules(api)
	if err != nil {
		return err
	}

	var filteredIpRules []*armsql.FirewallRule
	for _, ipRule := range ipRules {
		// make sure ipranges of a single ip are also removed
		if strings.TrimSuffix(*ipRule.Properties.StartIPAddress, "/32") != ip {
			filteredIpRules = append(filteredIpRules, ipRule)
		}
	}

	_, err = client.Replace(
		context.Background(),
		*rg,
		*s.Name,
		armsql.FirewallRuleList{Values: filteredIpRules},
		&armsql.FirewallRulesClientReplaceOptions{},
	)

	return err
}

func (s *PrickableSQLServer) ListIPRules(api *prick.Api) ([]*armsql.FirewallRule, error) {
	client := api.SqlFirewall
	rg, err := common.ExtractResourceGroup(s.ID)
	if err != nil {
		return nil, err
	}

	var fireWallRules []*armsql.FirewallRule

	fireWallRulePager := client.NewListByServerPager(*rg, *s.Name, &armsql.FirewallRulesClientListByServerOptions{})
	for fireWallRulePager.More() {
		res, err := fireWallRulePager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		rules := res.FirewallRuleListResult.Value
		fireWallRules = append(fireWallRules, rules...)
	}

	return fireWallRules, nil
}

func (s *PrickableSQLServer) ListPokes(api *prick.Api) ([]*common.Poke, error) {
	fireWallRules, err := s.ListIPRules(api)
	if err != nil {
		return nil, err
	}

	var pokes []*common.Poke
	for _, rule := range fireWallRules {
		pokes = append(
			pokes,
			&common.Poke{
				Name:           *rule.Name,
				StartIpAddress: *rule.Properties.StartIPAddress,
				EndIpAddress:   *rule.Properties.EndIPAddress,
			},
		)
	}

	return pokes, nil
}
