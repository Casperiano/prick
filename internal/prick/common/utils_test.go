package common

import (
	"testing"
)

func TestResourceGroupExtraction(t *testing.T) {
	id := "/subscriptions/295b5722-2faa-416b-b22e-e5926000eefd/resourceGroups/rg-adp-adp/providers/Microsoft.Storage/storageAccounts/stadpadpdev"
	exRg, err := ExtractResourceGroup(&id)

	if *exRg != "rg-adp-adp" || err != nil {
		t.Fatalf(`Extracted %q, %v, nil`, *exRg, err)
	}
}
