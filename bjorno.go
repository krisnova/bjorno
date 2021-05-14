package bjorno

import (
	"net/http"

	"github.com/kris-nova/logger"
)

type ServerConfig struct {
	LogVerbosity   int    // 3
	ServeDirectory string // /
	BindAddress    string // localhost:80

	// InterpolateExtensions are the names
	// of the interpolation file extensions
	// to parse.
	//
	// Usually this is something like ".html"
	InterpolateExtensions []string

	// DefaultIndexFiles are the names of file
	// to look for if a directory is passed.
	//
	// Usually this is something like "index.html".
	DefaultIndexFiles []string

	// UseDefaultRootHandler gives us an easy way
	// to run the server in a more "concrete" and
	// reliable way while we are in alpha dev stages.
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
	EndpointRoot     string = "/"
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
		logger.BitwiseLevel = logger.LogSuccess | logger.LogAlways | logger.LogCritical | logger.LogWarning | logger.LogInfo
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
	logger.Info("Default Index Files: %s", cfg.DefaultIndexFiles)

	// Root (/) handler
	if cfg.UseDefaultRootHandler {
		logger.Info("Using default root handler")
		http.Handle(EndpointRoot, http.FileServer(http.Dir(cfg.ServeDirectory)))
	} else {
		logger.Info("Using bjorno root handler")
		http.Handle(EndpointRoot, NewRootHandler(cfg))
	}

	// Client (/client) handler
	http.Handle("/client", NewClientHandler(cfg))

	// Because we define custom handlers above we do not need to
	// pass in a "generic" handler here.
	return http.ListenAndServe(cfg.BindAddress, nil)
}
