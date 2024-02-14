package config

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("ENV", "local")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("AUTO_MIGRATIONS", true)
}

type Config struct {
	Env            string
	Port           string
	AutoMigrations bool
}

func NewConfig() *Config {
	return &Config{
		Env:            viper.GetString("ENV"),
		Port:           viper.GetString("PORT"),
		AutoMigrations: viper.GetBool("AUTO_MIGRATIONS"),
	}
}
