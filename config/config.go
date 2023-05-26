package config

import "go.uber.org/zap"

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
	TLS  *TLSSettings       `mapstructure:"tls"`
}

type HTTPServerSettings struct {
	Endpoint string `mapstructure:"endpoint"`
}

type TLSSettings struct {
	CertFile    string `mapstructure:"cert_file"`
	CertKeyFile string `mapstructure:"cert_key_file"`
	Insecure    bool   `mapstructure:"insecure"`
}

func New() *Cfg {
	cfg := &Cfg{
		Logger: zap.Must(zap.NewDevelopment()),
		Server: ServerSettings{
			TLS: &TLSSettings{},
		},
	}

	return cfg
}
