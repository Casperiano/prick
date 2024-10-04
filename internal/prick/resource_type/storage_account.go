package resource_type

import (
	"context"
	prick "prick/internal/prick"
	"prick/internal/prick/common"
	"prick/internal/prick/interfaces"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
)

var _ interfaces.Prickable = &PrickableStorageAccount{}

type PrickableStorageAccount armstorage.Account

func (psa *PrickableStorageAccount) Poke(api *prick.Api) error {
	client := api.StorageAccount
	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}

	ipRules, err := psa.ListIPRules(api)

	if err != nil {
		return err
	}

	rg, err := common.ExtractResourceGroup(psa.ID)
	if err != nil {
		return err
	}

	updatedIpRules := append(ipRules, &armstorage.IPRule{IPAddressOrRange: &ip})

	_, err = client.Update(context.Background(), *rg, *psa.Name,
		armstorage.AccountUpdateParameters{
			Properties: &armstorage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: &armstorage.NetworkRuleSet{
					IPRules: updatedIpRules,
				},
			},
		},
		&armstorage.AccountsClientUpdateOptions{})
	return err
}

func (psa *PrickableStorageAccount) Patch(api *prick.Api) error {
	client := api.StorageAccount
	rg, err := common.ExtractResourceGroup(psa.ID)
	if err != nil {
		return err
	}

	ip, err := common.GetIPAddress()
	if err != nil {
		return err
	}
	ipRules, err := psa.ListIPRules(api)
	if err != nil {
		return err
	}
	var filteredIpRules []*armstorage.IPRule
	for _, ipRule := range ipRules {
		// make sure ipranges of a single ip are also removed
		if strings.TrimSuffix(*ipRule.IPAddressOrRange, "/32") != ip {
			filteredIpRules = append(filteredIpRules, ipRule)
		}
	}

	_, err = client.Update(context.Background(), *rg, *psa.Name,
		armstorage.AccountUpdateParameters{
			Properties: &armstorage.AccountPropertiesUpdateParameters{
				NetworkRuleSet: &armstorage.NetworkRuleSet{
					IPRules: filteredIpRules,
				},
			},
		},
		&armstorage.AccountsClientUpdateOptions{})
	return err
}

func ListStorageAccounts(api *prick.Api, options *ListStorageAccountsOptions) ([]*PrickableStorageAccount, error) {
	client := api.StorageAccount
	// TODO: refactor with generics to avoid repetition
	var accounts []*PrickableStorageAccount
	switch options.ResourceGroup {
	case "":
		pager := client.NewListPager(&armstorage.AccountsClientListOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, sa := range res.Value {
				newSa := PrickableStorageAccount(*sa)
				accounts = append(accounts, &newSa)
			}
		}
	default:
		pager := client.NewListByResourceGroupPager(options.ResourceGroup, &armstorage.AccountsClientListByResourceGroupOptions{})
		for pager.More() {
			res, err := pager.NextPage(context.Background())
			if err != nil {
				return nil, err
			}

			for _, sa := range res.Value {
				newSa := PrickableStorageAccount(*sa)
				accounts = append(accounts, &newSa)
			}
		}
	}

	return accounts, nil
}

type ListStorageAccountsOptions struct {
	ResourceGroup string
}

func (psa *PrickableStorageAccount) GetName() string {
	return *psa.Name
}

func (psa *PrickableStorageAccount) GetLocation() string {
	return *psa.Location
}

func (psa *PrickableStorageAccount) GetType() common.ResourceType {
	return common.ResourceTypeStorageAccount
}

func (psa *PrickableStorageAccount) ListPokes(api *prick.Api) ([]*common.Poke, error) {
	rules, err := psa.ListIPRules(api)
	if err != nil {
		return nil, err
	}

	pokes := []*common.Poke{}
	for _, rule := range rules {
		startIp, endIp, err := common.ParseCidr(*rule.IPAddressOrRange)
		if err != nil {
			return nil, err
		}
		pokes = append(pokes, &common.Poke{StartIpAddress: startIp, EndIpAddress: endIp})
	}
	return pokes, nil
}

func (psa *PrickableStorageAccount) ListIPRules(api *prick.Api) ([]*armstorage.IPRule, error) {
	client := api.StorageAccount
	rg, err := common.ExtractResourceGroup(psa.ID)
	if err != nil {
		return nil, err
	}
	props, err := client.GetProperties(context.Background(), *rg, *psa.Name, &armstorage.AccountsClientGetPropertiesOptions{})
	if err != nil {
		return nil, err
	}
	return props.Properties.NetworkRuleSet.IPRules, nil

}
