package main

import (
	"fmt"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration from a file
func InitConfig(configPath string) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("configuration")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading configuration: %s \n", err))
	}
	if debug {
		fmt.Println("> Configuration loaded from file:", viper.ConfigFileUsed())
	}
}
