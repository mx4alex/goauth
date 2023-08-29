package config

import (
	"github.com/spf13/viper"
)

const (
	ConfigFilePath = "./"
	ConfigFileName = "config"
)

type Config struct {
	HostAddr string       `mapstructure:"host_addr"`
	Auth  	 AuthConfig	  `mapstructure:"auth"`
	MongoDB  MongoConfig  `mapstructure:"mongo"`
}

type AuthConfig struct {
	SigningKey string  `mapstructure:"signing_key"`
	HashSalt   string  `mapstructure:"hash_salt"`
	TokenTTL   int     `mapstructure:"token_ttl"`
}

type MongoConfig struct {
	Url 	   string  `mapstructure:"url"`
	Name 	   string  `mapstructure:"name"`
	Collection string  `mapstructure:"collection"`
}

func New() (Config, error) {
	vpr := viper.New()
	vpr.AddConfigPath(ConfigFilePath)
	vpr.SetConfigName(ConfigFileName)

	if err := vpr.ReadInConfig(); err != nil {
		return Config{}, err
	}

	var result Config
	if err := vpr.Unmarshal(&result); err != nil {
		return Config{}, err
	}

	return result, nil
}
