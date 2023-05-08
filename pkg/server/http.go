package server

import (
	"github.com/dose-na-nuvem/customers/config"
)

func Serve(cfg *config.Cfg) {
	cfg.Logger.Info("serving...\n")
}
