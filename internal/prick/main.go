package prick

import (
	"prick/internal/prick/common"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
)

type Api struct {
	ResourceGroup    *armresources.ResourceGroupsClient
	KeyVault         *armkeyvault.VaultsClient
	StorageAccount   *armstorage.AccountsClient
	SqlServer        *armsql.ServersClient
	SqlFirewall      *armsql.FirewallRulesClient
	SynapseWorkspace *armsynapse.WorkspacesClient
	SynapseFirewall  *armsynapse.IPFirewallRulesClient
}

func NewApi() (*Api, error) {
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}
	sid, err := common.GetSubscriptionId()
	if err != nil {
		return nil, err
	}

	rgClient, err := NewResourceGroupClient(credential, sid)
	if err != nil {
		return nil, err
	}

	kvClient, err := NewKeyVaultClient(credential, sid)
	if err != nil {
		return nil, err
	}

	saClient, err := NewAccountsClient(credential, sid)
	if err != nil {
		return nil, err
	}

	ssClient, err := NewSqlServerClient(credential, sid)
	if err != nil {
		return nil, err
	}

	fwClient, err := NewSqlServerFirewallClient(credential, sid)
	if err != nil {
		return nil, err
	}

	swClient, err := NewSynapseClient(credential, sid)
	if err != nil {
		return nil, err
	}

	sfwClient, err := NewSynapseFirewallClient(credential, sid)
	if err != nil {
		return nil, err
	}

	return &Api{rgClient, kvClient, saClient, ssClient, fwClient, swClient, sfwClient}, nil
}

func NewResourceGroupClient(credential *azidentity.DefaultAzureCredential, sid string) (*armresources.ResourceGroupsClient, error) {
	resourceGroupClient, err := armresources.NewResourceGroupsClient(sid, credential, nil)
	if err != nil {
		return nil, err
	}
	return resourceGroupClient, nil
}

func NewKeyVaultClient(credential *azidentity.DefaultAzureCredential, sid string) (*armkeyvault.VaultsClient, error) {
	factory, err := armkeyvault.NewClientFactory(sid, credential, &policy.ClientOptions{})
	if err != nil {
		return nil, err
	}

	return factory.NewVaultsClient(), nil
}

func NewAccountsClient(credential *azidentity.DefaultAzureCredential, sid string) (*armstorage.AccountsClient, error) {
	factory, err := armstorage.NewClientFactory(sid, credential, &policy.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return factory.NewAccountsClient(), nil
}

func NewSqlServerClient(credential *azidentity.DefaultAzureCredential, sid string) (*armsql.ServersClient, error) {
	factory, err := armsql.NewClientFactory(sid, credential, &policy.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return factory.NewServersClient(), nil
}

func NewSqlServerFirewallClient(credential *azidentity.DefaultAzureCredential, sid string) (*armsql.FirewallRulesClient, error) {
	factory, err := armsql.NewClientFactory(sid, credential, &policy.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return factory.NewFirewallRulesClient(), nil
}

func NewSynapseClient(credential *azidentity.DefaultAzureCredential, sid string) (*armsynapse.WorkspacesClient, error) {
	factory, err := armsynapse.NewClientFactory(sid, credential, &policy.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return factory.NewWorkspacesClient(), nil
}

func NewSynapseFirewallClient(credential *azidentity.DefaultAzureCredential, sid string) (*armsynapse.IPFirewallRulesClient, error) {
	factory, err := armsynapse.NewClientFactory(sid, credential, &policy.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return factory.NewIPFirewallRulesClient(), nil
}
