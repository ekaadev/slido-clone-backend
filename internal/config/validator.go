package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// NewValidator function untuk create validator yang digunakan untuk validasi struct
func NewValidator(viper *viper.Viper) *validator.Validate {
	return validator.New()
}
