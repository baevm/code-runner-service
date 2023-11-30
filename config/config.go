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
		MaxCount int           `mapstructure:"max_count"`
		MaxTime  time.Duration `mapstructure:"max_time"`
	} `mapstructure:"containers"`

	RabbitMQ struct {
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		VHost    string `mapstructure:"vhost"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"rabbitmq"`
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
