package bjorno

import (
	"testing"

	"github.com/kris-nova/logger"
)

func TestRunServerWithConfig(t *testing.T) {
	cfg := &ServerConfig{
		ServeDirectory: "/tmp",
		BindAddress:    ":1314",
		LogVerbosity:   1,
	}
	logger.Always("Running test server: %s", cfg.BindAddress)
	// TODO we don't really have anything to test yet
	// this is all in the Go standard library
}
