package config

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reflect"
)

const (
	configFile = "app.env"
)

type Config struct {
	DBConfig      DBConfig    `mapstructure:"DATABASE"`
	ServerAddress string      `mapstructure:"SERVER_ADDRESS"`
	CacheConfig   CacheConfig `mapstructure:"CACHE"`
	AuthConfig    AuthConfig  `mapstructure:"AUTH"`
}

var config *Config

func GetConfig(paths ...string) (*Config, error) {
	if config != nil {
		return config, nil
	}
	conf, err := LoadConfig(paths...)
	config = conf
	return conf, err
}

func LoadConfig(paths ...string) (conf *Config, err error) {
	v := viper.New()

	v.SetConfigType("yaml")
	if err = v.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
		logrus.Fatalf("error loading default config -> %s", err.Error())
	}

	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	for _, path := range paths {
		v.AddConfigPath(path)
	}
	v.AutomaticEnv()

	switch err = v.MergeInConfig(); err.(type) {
	case nil:
	case viper.ConfigFileNotFoundError:
		logrus.Warnf("config path error  -> %s .\n using default", err.Error())
	default:
		logrus.Infof("error type : %s", reflect.TypeOf(err))
		logrus.Warnf("error loading conf -> %s", err.Error())
	}

	if err = v.UnmarshalExact(&conf); err != nil {
		logrus.Fatalf("couldn't unmarshal conf -> %s", err.Error())
	}
	return
}
