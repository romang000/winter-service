package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServiceConfig struct {
	ServerPort string `mapstructure:"server_port"`
	Address    string `mapstructure:"address"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type AppConfig struct {
	Service  ServiceConfig  `mapstructure:"service"`
	Database DatabaseConfig `mapstructure:"database"`
}

func LoadConfig(path string) (AppConfig, error) {
	var config AppConfig
	
	viper.SetConfigFile(path)
	
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}
	
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return config, nil
}
