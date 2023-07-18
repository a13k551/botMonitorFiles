package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Load() Config {

	viper.SetConfigName("conf")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../..")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var Config Config

	err = mapstructure.Decode(viper.AllSettings(), &Config)

	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return Config
}

type Config struct {
	Token       string
	Chatid      int
	Path        string
	Mask        string
	MinDate     string
	UserDB      string
	PassDB      string
	NameDB      string
	TableNamedb string
	FieldnameDB string
	Port        uint16
}
