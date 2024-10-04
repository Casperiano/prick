//go:build integration

package resource_type

import (
	prick "prick/internal/prick"
	"testing"
)

func TestListKeyVaults(t *testing.T) {
	api, err := prick.NewApi()
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	kv, err := ListKeyVaults(api, &ListKeyVaultsOptions{})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("KeyVaults: %v", *kv[0].Name)
}
