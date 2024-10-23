package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Address struct {
		Address string `mapstructure:"address"`
	}

	StressTest struct {
		AddResponsesCount    int `mapstructure:"add_responses_count"`
		GiveResponsesCount   int `mapstructure:"give_responses_count"`
		RefundResponsesCount int `mapstructure:"refund_responses_count"`
		ReturnResponsesCount int `mapstructure:"return_responses_count"`
	}

	Config struct {
		GRPC Address    `mapstructure:"grpc"`
		Test StressTest `mapstructure:"test"`
	}
)

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("stress_test")
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
