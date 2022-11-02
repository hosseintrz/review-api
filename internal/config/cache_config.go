package config

import "time"

type CacheConfig struct {
	Type            string        `mapstructure:"TYPE"`
	Address         string        `mapstructure:"ADDRESS"`
	Username        string        `mapstructure:"USERNAME"`
	PASSWORD        string        `mapstructure:"PASSWORD"`
	DB              int           `mapstructure:"DB"`
	ConnMaxIdleTime time.Duration `mapstructure:"CONN_MAX_IDLE_TIME"`
}
