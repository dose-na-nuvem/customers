package config

import "go.uber.org/zap"

type Cfg struct {
	DBType string `mapstruct:"db_type"`

	Logger *zap.Logger
}

func New() *Cfg {
	cfg := &Cfg{
		Logger: zap.Must(zap.NewDevelopment()),
	}

	return cfg
}
