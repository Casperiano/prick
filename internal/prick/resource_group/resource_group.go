package resource_group

import (
	"context"
	prick "prick/internal/prick"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

func ListResourceGroups(api *prick.Api) ([]*armresources.ResourceGroup, error) {
	client := api.ResourceGroup

	var resourceGroups []*armresources.ResourceGroup
	pager := client.NewListPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		resourceGroups = append(resourceGroups, page.ResourceGroupListResult.Value...)
	}

	return resourceGroups, nil
}
