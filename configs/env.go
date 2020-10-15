package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	RepositoryJsonPath string
	Port               int
	Files              []string
}

func New() *Config {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetDefault("PORT", 8080)

	return &Config{
		RepositoryJsonPath: viper.GetString("JSON_FILES_PATH"),
		Port:               viper.GetInt("PORT"),
		Files:              viper.GetStringSlice("JSON_FILES"),
	}
}
