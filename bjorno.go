package bjorno

import (
	"net/http"

	"github.com/kris-nova/logger"
)

type ServerConfig struct {
	LogVerbosity   int    // 3
	ServeDirectory string // /
	BindAddress    string // localhost:80
}

func RunServer(cfg *ServerConfig) error {
	if cfg.LogVerbosity == 1 {
		logger.BitwiseLevel = logger.LogEverything
	} else {
		logger.BitwiseLevel = logger.LogCritical | logger.LogWarning
	}
	logger.Info("ServeDirectory: %s", cfg.ServeDirectory)
	logger.Info("BindAddress: %s", cfg.BindAddress)
	logger.Info("Verbosity: %d", cfg.LogVerbosity)
	fileServer := http.FileServer(http.Dir(cfg.ServeDirectory))
	http.Handle("/", fileServer)
	return http.ListenAndServe(cfg.BindAddress, nil)
}
