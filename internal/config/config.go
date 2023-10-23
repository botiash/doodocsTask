package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

type (
	Config struct {
		Smtp_server string `mapstructure:"smtp_server"`
		Username    string `mapstructure:"username"`
		Password    string `mapstructure:"password"`
	}
)

var (
	instance Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		viper.SetConfigFile("config.yaml")
		// Load .env file
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config.yaml file: %v", err)
		}

		// Unmarshal configuration from Viper
		if err := viper.Unmarshal(&instance); err != nil {
			log.Fatalf("Error unmarshaling configuration: %v", err)
		}

		fmt.Println("instance of the application: ", instance)

	})

	return &instance
}
