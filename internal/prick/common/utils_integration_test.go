//go:build integration

package common

import (
	"testing"

	"github.com/google/uuid"
)

func TestGetIPAddress(t *testing.T) {
	ip := GetIPAddress()
	if ip == "" {
		t.Fatalf("IP address is empty")
	}
}

func TestGetSubscriptionId(t *testing.T) {
	id, err := GetSubscriptionId()
	if err != nil {
		t.Fatal("Error is nil")
	}

	_, err = uuid.Parse(id)
	if err != nil {
		t.Fatal("Not a valid uuid")
	}
}
