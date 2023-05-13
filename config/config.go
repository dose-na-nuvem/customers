package config

import "go.uber.org/zap"

type Cfg struct {
	Database Database `mapstructure:"db"`

	Logger *zap.Logger
}

type Database struct {
	Type     string `mapstructure:"type"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func New() *Cfg {
	cfg := &Cfg{
		Logger: zap.Must(zap.NewDevelopment()),
	}

	return cfg
}
