package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func New() *viper.Viper {
	cfg := viper.New()

	// Servers part of config
	cfg.SetDefault("LISTEN", ":8000")

	// Data for connecting to the database
	cfg.SetDefault("DB_HOST", "localhost")
	cfg.SetDefault("DB_PORT", "5432")
	cfg.SetDefault("DB_USER", "postgresql")
	cfg.SetDefault("DB_PASSWORD", "postgresql")
	cfg.SetDefault("DB_NAME", "postgresql")


	cfg.AutomaticEnv()

	cfg.Set("DB", dbConnectionString(cfg.GetString("DB_HOST"),
		cfg.GetString("DB_PORT"),
		cfg.GetString("DB_USER"),
		cfg.GetString("DB_PASSWORD"),
		cfg.GetString("DB_NAME")))

	return cfg
}
func dbConnectionString(host string, port string, user string, password string, dbname string) string {
	return fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}


