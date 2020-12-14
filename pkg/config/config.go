package config

import (
	"github.com/spf13/viper"
)

func New() *viper.Viper {
	cfg := viper.New()

	// Servers part of config
	cfg.SetDefault("LISTEN", ":8000")

	cfg.AutomaticEnv()

	return cfg
}


