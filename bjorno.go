package bjorno

import (
	"github.com/kris-nova/logger"
	"net/http"
)

type ServerConfig struct {
	LogVerbosity   int    // 3
	ServeDirectory string // /
	BindAddress    string // localhost:80

	UseDefaultRootHandler bool

	// The Server can have custom global response content
	// for specific HTTP Error codes
	// These can be defined at runtime.
	Content404 []byte
	Content500 []byte
	Content5XX []byte
}

const (
	StatusDefault404 string = `404 not found (bjorno)`
	StatusDefault500 string = `500 server error (bjorno)`
	StatusDefault5XX string = `5xx server error (bjorno)`
	EndpointRoot string = "/"
)

// RunServer should be stateless by design
// We never want a server.DoThing() paradigm
// It either enters via config or we need to
// figure something else out.
func RunServer(cfg *ServerConfig) error {
	v := cfg.LogVerbosity
	switch v {
	case 4:
		logger.BitwiseLevel = logger.LogEverything
	case 3:
		logger.BitwiseLevel = logger.LogSuccess | logger.LogAlways | logger.LogCritical | logger.LogWarning  | logger.LogInfo
	case 2:
		logger.BitwiseLevel = logger.LogSuccess | logger.LogAlways | logger.LogCritical | logger.LogWarning
	case 1:
		logger.BitwiseLevel = logger.LogSuccess | logger.LogAlways | logger.LogCritical
	case 0:
		logger.BitwiseLevel = logger.LogAlways
	}

	logger.Info("ServeDirectory: %s", cfg.ServeDirectory)
	logger.Info("BindAddress: %s", cfg.BindAddress)
	logger.Info("Verbosity: %d", cfg.LogVerbosity)
	logger.Info("Log Level: %d", logger.BitwiseLevel)


	// We can have different endpoints that perform different
	// ways.
	//
	//http.Handle("/api", apiHandler)
	//http.Handle("/example", exampleHandler)
	if cfg.UseDefaultRootHandler {
		logger.Info("Using default root handler")
		http.Handle(EndpointRoot, http.FileServer(http.Dir(cfg.ServeDirectory)))
	}else {
		logger.Info("Using bjorno root handler")
		http.Handle(EndpointRoot, NewRootHandler(cfg))
	}


	// Because we define custom handlers above we do not need to
	// pass in a "generic" handler here.
	return http.ListenAndServe(cfg.BindAddress, nil)
}



