package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	PORT        int    `mapstructure:"PORT"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     int    `mapstructure:"DB_PORT"`
	DB_NAME     string `mapstructure:"DB_NAME"`
	JWT_SECRET  string `mapstructure:"JWT_SECRET"`
}

func New() (*Config, error) {
	initLogging("info")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Error("config file not found")
		return nil, fmt.Errorf("config file not found")
	} else {
		zap.S().Info("config file found")
	}
	c := &Config{}

	if err := viper.Unmarshal(c); err != nil {
		zap.S().Error("unable to decode into struct", zap.Error(err))
		return nil, fmt.Errorf("unable to decode into struct")
	} else {
		zap.S().Info("config file loaded")
	}

	return c, nil
}
