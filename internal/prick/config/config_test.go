package config

import (
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestParseConfig2(t *testing.T) {
	rawConfig := `
pricks:
  fab:
    storage_accounts:
      stadpfabprd:
        - startIp: "1.2.3.4"
          endIp: "4.3.2.1"
`

	viper.SetConfigType("yaml")
	err := viper.ReadConfig(strings.NewReader(rawConfig))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var config Config
	if err = viper.Unmarshal(&config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Equal(t, config.Pricks["fab"]["storage_accounts"]["stadpfabprd"][0].StartIp, "1.2.3.4")
	assert.Equal(t, config.Pricks["fab"]["storage_accounts"]["stadpfabprd"][0].EndIp, "4.3.2.1")
}
