package config

import (
	"go.uber.org/zap"
	"time"
)

type Cfg struct {
	Database Database       `mapstructure:"db"`
	Server   ServerSettings `mapstructure:"server"`

	Logger *zap.Logger
}

type Database struct {
	Type     string `mapstructure:"type"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type ServerSettings struct {
	HTTP HTTPServerSettings `mapstructure:"http"`
}

type HTTPServerSettings struct {
	Endpoint          string        `mapstructure:"endpoint"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
}

func New() *Cfg {
	cfg := &Cfg{
		Logger: zap.Must(zap.NewDevelopment()),
	}

	return cfg
}
