//go:build integration

package resource_type

import (
	prick "prick/internal/prick"
	"testing"
)

func TestListStorageAccounts(t *testing.T) {
	api, err := prick.NewApi()
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	accounts, err := ListStorageAccounts(api, &ListStorageAccountsOptions{ResourceGroup: ""})
	if err != nil {
		t.Fatalf("Error:  %v", err)
	}
	t.Log(*accounts[0].Name)
}
