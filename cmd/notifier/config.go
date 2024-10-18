package main

import (
	"fmt"

	"github.com/spf13/viper"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka"
)

type (
	Kafka struct {
		Config  kafka.Config `mapstructure:"config"`
		Topics  []string     `mapstructure:"topics"`
		GroupID string       `mapstructure:"group_id"`
	}

	Config struct {
		Kafka Kafka `mapstructure:"kafka"`
	}
)

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("notifier")
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
