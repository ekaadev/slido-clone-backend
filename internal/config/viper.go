package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// NewViper function untuk load file configuration
// Example: load .env and config
func NewViper() *viper.Viper {
	config := viper.New()

	// read configuration from config.json (general configuration)
	config.SetConfigFile("config.json")
	if err := config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// read configuration from .env (sensitive/secret configuration)
	// If .env is absent (e.g., running in Docker), fall back to OS environment variables
	// loaded via AutomaticEnv() above.
	config.SetConfigFile(".env")
	config.SetConfigType("env")
	config.AutomaticEnv()
	if err := config.MergeInConfig(); err != nil {
		if !os.IsNotExist(err) {
			panic(fmt.Errorf("Fatal error env file: %s \n", err))
		}
	}

	return config
}
