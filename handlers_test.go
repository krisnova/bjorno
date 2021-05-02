package bjorno

import (
	"testing"
)

func TestRootExample(t *testing.T) {
	cfg := &ServerConfig{
		ServeDirectory: "testWebsite",
		BindAddress:    ":1314",
		LogVerbosity:   1,
	}
	err := RunServer(cfg)
	if err != nil {
		t.Errorf("Error running server: %v", err)
	}
}