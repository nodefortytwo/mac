package config

import (
	"github.com/spf13/viper"
)

type Config struct {
}

func New() Config {
	return Config{}
}

func (c Config) GetCodeRoot() string {
	viper.SetDefault("CodeRoot", "/var/www")
	return viper.GetString("CodeRoot")
}

func (c Config) GetQuitExclusions() []string {
	viper.SetDefault("QuitExclusions", []string{"terminal"})
	return viper.GetStringSlice("QuitExclusions")
}
