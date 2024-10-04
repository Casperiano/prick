package resource_type

import (
	"context"
	"fmt"
	prick "prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
)

var _ interfaces.Prickable = &PrickableKeyVault{}

type PrickableKeyVault armkeyvault.Vault

func ListKeyVaults(api *prick.Api, options *ListKeyVaultsOptions) ([]*PrickableKeyVault, error) {
	client := api.KeyVault

	var vaults []*PrickableKeyVault
	switch options.ResourceGroup {
	case "":
		pager := client.NewListBySubscriptionPager(&armkeyvault.VaultsClientListBySubscriptionOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, kv := range res.Value {
				newKv := PrickableKeyVault(*kv)
				vaults = append(vaults, &newKv)
			}
		}
	default:
		pager := client.NewListByResourceGroupPager(options.ResourceGroup, &armkeyvault.VaultsClientListByResourceGroupOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, kv := range res.Value {
				newKv := PrickableKeyVault(*kv)
				vaults = append(vaults, &newKv)
			}
		}
	}

	return vaults, nil
}

type ListKeyVaultsOptions struct{ ResourceGroup string }

func (kv *PrickableKeyVault) Poke(api *prick.Api) error {
	client := api.KeyVault
	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}

	rg, err := common.ExtractResourceGroup(kv.ID)
	if err != nil {
		return err
	}

	ipRules, err := kv.ListIPRules(api)
	if err != nil {
		return err
	}

	ipInIpRules, err := common.IpInIpRules(ip, ipRules)
	if err != nil {
		return err
	}
	if ipInIpRules {
		return fmt.Errorf("IP is already in the IP rules: %v", err)
	}

	updatedIpRules := append(ipRules, &armkeyvault.IPRule{Value: &ip})

	// In order for ACLs to apply, the key vault should have been created
	// with the "Allow public access from specific virtual networks and IP addresses" configuration
	// see https://github.com/hashicorp/terraform-provider-azurerm/issues/25414
	pna := string(armkeyvault.PublicNetworkAccessEnabled)
	_, err = client.Update(context.Background(), *rg, *kv.Name,
		armkeyvault.VaultPatchParameters{
			Properties: &armkeyvault.VaultPatchProperties{
				NetworkACLs: &armkeyvault.NetworkRuleSet{
					IPRules: updatedIpRules,
				},
				PublicNetworkAccess: &pna,
			},
		},
		&armkeyvault.VaultsClientUpdateOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (kv *PrickableKeyVault) Patch(api *prick.Api) error {
	client := api.KeyVault
	rg, err := common.ExtractResourceGroup(kv.ID)
	if err != nil {
		return err
	}

	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}
	ipRules, err := kv.ListIPRules(api)
	if err != nil {
		return err
	}
	var filteredIpRules []*armkeyvault.IPRule
	for _, ipRule := range ipRules {
		if *ipRule.Value != fmt.Sprintf("%s/32", ip) { // Azure adds /32 ietself when poking with a single IP
			filteredIpRules = append(filteredIpRules, ipRule)
		}
	}

	_, err = client.Update(context.Background(), *rg, *kv.Name,
		armkeyvault.VaultPatchParameters{
			Properties: &armkeyvault.VaultPatchProperties{
				NetworkACLs: &armkeyvault.NetworkRuleSet{IPRules: filteredIpRules},
			},
		},
		&armkeyvault.VaultsClientUpdateOptions{})
	return err
}

func (kv *PrickableKeyVault) GetName() string {
	return *kv.Name
}

func (kv *PrickableKeyVault) GetLocation() string {
	return *kv.Location
}

func (kv *PrickableKeyVault) GetType() common.ResourceType {
	return common.ResourceTypeKeyVault
}

func (kv *PrickableKeyVault) ListPokes(api *prick.Api) ([]*common.Poke, error) {
	ipRules, err := kv.ListIPRules(api)
	if err != nil {
		return nil, err
	}

	pokes := []*common.Poke{}
	for _, ipRule := range ipRules {
		startIp, endIp, err := common.ParseCidr(*ipRule.Value)
		if err != nil {
			return nil, err

		}
		pokes = append(pokes, &common.Poke{StartIpAddress: startIp, EndIpAddress: endIp})
	}

	return pokes, nil
}

func (kv *PrickableKeyVault) ListIPRules(api *prick.Api) ([]*armkeyvault.IPRule, error) {
	client := api.KeyVault

	rg, err := common.ExtractResourceGroup(kv.ID)
	if err != nil {
		return nil, err
	}

	props, err := client.Get(context.Background(), *rg, *kv.Name, &armkeyvault.VaultsClientGetOptions{})
	if err != nil {
		return nil, err
	}

	if props.Properties.NetworkACLs == nil {
		return nil, nil
	}

	return props.Properties.NetworkACLs.IPRules, nil
}
