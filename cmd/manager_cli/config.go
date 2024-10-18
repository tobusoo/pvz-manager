package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Address struct {
		Address string `mapstructure:"address"`
	}

	Config struct {
		GRPC Address `mapstructure:"grpc"`
	}
)

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("manager_cli")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("LoadConfig ReadInConfig: %w", err)
	}

	c := &Config{}
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("LoadConfig Unmarshal: %w", err)
	}

	return c, nil
}
