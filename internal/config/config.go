package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func GetConf() map[string]interface{} {

	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	conf := viper.AllSettings()

	return conf
}
