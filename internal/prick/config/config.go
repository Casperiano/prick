package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Pricks map[string]ResourceType `mapstructure:"pricks"`
}

type ResourceType map[string]Resources

type Resources map[string][]IpRule

type IpRule struct {
	StartIp string `mapstructure:"startIp"`
	EndIp   string `mapstructure:"endIp"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName(".prick")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
