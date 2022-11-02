package config

type DBConfig struct {
	Driver string `mapstructure:"DRIVER"`
	Source string `mapstructure:"SOURCE"`
}
