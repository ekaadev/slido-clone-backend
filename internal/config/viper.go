package config

import (
	"errors"
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
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// read configuration from .env (sensitive/secret configuration)
	// If .env is absent (e.g., running in Docker), fall back to OS environment variables
	// loaded via AutomaticEnv() above.
	config.SetConfigFile(".env")
	config.SetConfigType("env")
	config.AutomaticEnv()
	if err := config.MergeInConfig(); err != nil {
		// Handle both viper.ConfigFileNotFoundError (config name search) and
		// os.PathError/syscall.ENOENT (explicit file path not found).
		// If .env is absent, fall back to OS environment variables via AutomaticEnv().
		var notFound viper.ConfigFileNotFoundError
		if !os.IsNotExist(err) && !errors.As(err, &notFound) {
			panic(fmt.Errorf("fatal error env file: %s", err))
		}
	}

	return config
}
