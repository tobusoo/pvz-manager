package main

import (
	"fmt"

	"github.com/spf13/viper"
	"gitlab.ozon.dev/chppppr/homework/internal/infra/kafka"
)

type (
	Address struct {
		Address string `mapstructure:"address"`
	}

	Kafka struct {
		Config kafka.Config `mapstructure:"config"`
		Topic  string       `mapstructure:"topic"`
	}

	Config struct {
		GRPC    Address `mapstructure:"grpc"`
		HTPP    Address `mapstructure:"http"`
		Swagger Address `mapstructure:"swagger"`
		Kafka   Kafka   `mapstructure:"kafka"`
	}
)

func LoadConfig() (*Config, error) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("manager_service")
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
