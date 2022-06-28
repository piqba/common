package config

import (
	"github.com/spf13/viper"
	"os"
	"strings"
)

// LoadConfig load the config file from `path`, `name` and filetype
func LoadConfig(path, name, filetype string) (*viper.Viper, error) {
	vp := viper.New()
	vp.SetConfigName(name)
	vp.SetConfigType(filetype)
	vp.AddConfigPath(path)

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, ErrConfigFileNotFound
		} else {
			return nil, err
		}
	}
	// logic for get ENV_VARS
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vp.SetEnvPrefix("APP")
	vp.AutomaticEnv()

	return vp, nil
}

// LoadEnvOrFallback load the env vars, if not exists return a fallback...
func LoadEnvOrFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
