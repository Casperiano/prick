package common

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	"github.com/stretchr/testify/assert"
)

func TestParseCidrSingleIp(t *testing.T) {
	ip := "1.1.1.1"

	start, end, err := ParseCidr(ip)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, ip, start)
	assert.Equal(t, ip, end)
}

func TestParseCidrRange(t *testing.T) {
	ip := "1.2.3.0/28" // 1.2.3.4 -> 1.2.3.20

	start, end, err := ParseCidr(ip)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, start, "1.2.3.0")
	assert.Equal(t, end, "1.2.3.15")
}

func TestIpInIpRules(t *testing.T) {

	testCases := map[string]struct {
		ip      string
		ipRules []string
		result  bool
	}{
		"simple": {
			ip:      "1.1.1.1",
			ipRules: []string{"1.1.1.0/28"},
			result:  true,
		},
		"not in range": {
			ip:      "2.143.54.2",
			ipRules: []string{"1.1.1.0/28"},
			result:  false,
		},
		"range no overlap": {
			ip:      "1.1.1.0/28",
			ipRules: []string{"2.2.2.2/28"},
			result:  false,
		},
		"range overlap": {
			ip:      "1.1.1.0/28",
			ipRules: []string{"1.1.1.0/16"},
			result:  true,
		},
		"match ip range with single ip rules": {
			ip:      "1.1.1.0/28",
			ipRules: []string{"3.2.4.0"},
			result:  false,
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ipRules := make([]*armkeyvault.IPRule, len(test.ipRules))
			for i, ipRule := range test.ipRules {
				ipRules[i] = &armkeyvault.IPRule{Value: &ipRule}
			}

			result, err := IpInIpRules(test.ip, ipRules)

			assert.NoError(t, err)
			assert.Equal(t, test.result, result)
		})

	}
}
