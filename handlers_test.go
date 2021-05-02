package bjorno

import (
	"github.com/kris-nova/logger"
	"testing"
)

func TestRootExample(t *testing.T) {
	cfg := &ServerConfig{
		ServeDirectory: "testWebsite",
		BindAddress:    ":1314",
		LogVerbosity:   1,
	}

	logger.Always("Running test server: %s", cfg.BindAddress)
	// TODO we don't really have anything to test yet
	// this is all in the Go standard library

	//err := RunServer(cfg)
	//if err != nil {
	//	t.Errorf("Error running server: %v", err)
	//}
}