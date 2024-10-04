package common

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	netipx "go4.org/netipx"
)

type ResourceType string

const (
	ResourceTypeStorageAccount   ResourceType = "storage-account"
	ResourceTypeKeyVault         ResourceType = "keyvault"
	ResourceTypeSynapseWorkspace ResourceType = "synapse"
	ResourceTypeSQLServer        ResourceType = "sql-server"
)

func ResourceTypes() []ResourceType {
	return []ResourceType{
		ResourceTypeStorageAccount,
		ResourceTypeKeyVault,
		ResourceTypeSynapseWorkspace,
		ResourceTypeSQLServer,
	}
}

func (rt *ResourceType) String() string {
	return string(*rt)
}

func (rt *ResourceType) Set(v string) error {
	switch v {
	case "storage-account", "keyvault", "synapse", "sql-server":
		*rt = ResourceType(v)
		return nil
	default:
		return fmt.Errorf(`must be one of %v`, ResourceTypes())
	}
}

func (rt *ResourceType) Type() string {
	return "resourceType"
}

type Poke struct {
	Name           string
	StartIpAddress string
	EndIpAddress   string
}

func ParseCidr(cidrOrIp string) (string, string, error) {
	ip, err := netip.ParseAddr(cidrOrIp)

	if err != nil {
		ipNet, err := netip.ParsePrefix(cidrOrIp)
		if err != nil {
			return "", "", err
		}
		rangeIp := netipx.RangeOfPrefix(ipNet)

		return rangeIp.From().String(), rangeIp.To().String(), nil
	}

	return ip.String(), ip.String(), nil
}

func IpInIpRules(cidrOrIp string, ipRules []*armkeyvault.IPRule) (bool, error) {
	var ipNet *net.IPNet
	ip, err := netip.ParseAddr(cidrOrIp)
	if err != nil {
		_, ipNet, err = net.ParseCIDR(cidrOrIp)
		if err != nil {
			return false, err
		}
	}

	switch ip.IsValid() {
	case true:
		for _, ipRule := range ipRules {
			startIp, endIp, err := ParseCidr(*ipRule.Value)
			if err != nil {
				return false, err
			}
			if ip.String() <= endIp && ip.String() >= startIp {
				return true, nil
			}
		}
		return false, nil
	case false:
		for _, ipRule := range ipRules {
			_, err := netip.ParseAddr(*ipRule.Value)
			if err != nil {
				_, ruleNetIp, err := net.ParseCIDR(*ipRule.Value)
				if err != nil {
					return false, err
				}

				if ipNet.Contains(ruleNetIp.IP) || ruleNetIp.Contains(ipNet.IP) {
					return true, nil
				}
			}
		}
		return false, nil
	}
	// This should never be executed?
	return false, nil
}
