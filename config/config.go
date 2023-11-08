package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`

	Containers struct {
		MaxCount int8          `mapstructure:"max_count"`
		MaxTime  time.Duration `mapstructure:"max_time"`
	} `mapstructure:"containers"`
}

var config = new(Config)

func Load() error {
	viper.SetConfigType("toml")
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	return config
}
