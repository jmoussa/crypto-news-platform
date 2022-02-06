package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	configName  = "config"
	configPaths = []string{
		os.Getenv("CONFIG_LOCATION"),
	}
)

type Config struct {
	Twitter map[string]string `json:"twitter"`
}

func ParseConfig() Config {
	var config Config
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.SetConfigType("json")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(os.Getenv("CONFIG_LOCATION"))
	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
		return config
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Cannot unmarshal config %s", err)
		return config
	}
	return config
}
