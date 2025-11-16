package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewLogger function untuk create logrus yang digunakan untuk logging aplikasi
func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	// Set log level and format
	log.SetLevel(logrus.Level(viper.GetInt("log.level")))
	log.SetFormatter(&logrus.JSONFormatter{})

	return log
}
