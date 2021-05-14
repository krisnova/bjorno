package bjorno

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/kris-nova/logger"
)

var (
	TestConfig = &ServerConfig{
		ServeDirectory: "./bjorno.com",
		BindAddress:    "localhost:1313",
		LogVerbosity:   1,
		Content404:     []byte("oopsie"),
	}
)

func TestMain(m *testing.M) {
	go func() {
		err := Runtime(TestConfig, &EmptyProgram{})
		if err != nil {
			logger.Critical("Error starting test server: %v", err)
			os.Exit(1)
		}
	}()
	exitCode := m.Run()
	if exitCode != 0 {
		logger.Critical("Test Failure")
		os.Exit(exitCode)
	}
	os.Exit(0)
}

func Test200IndexHTML(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/index.html", TestConfig.BindAddress))
	if err != nil {
		t.Errorf("Unable to resolve index.html: %v", err)
	}
	if resp.ContentLength < 1 {
		t.Errorf("Missing index.html")
	}
}

func Test200Slash(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/", TestConfig.BindAddress))
	if err != nil {
		t.Errorf("Unable to resolve index.html: %v", err)
	}
	if resp.ContentLength < 1 {
		t.Errorf("Missing index.html")
	}
}

func Test200Empty(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s", TestConfig.BindAddress))
	if err != nil {
		t.Errorf("Unable to resolve index.html: %v", err)
	}
	if resp.ContentLength < 1 {
		t.Errorf("Missing index.html")
	}
}

func Test404Indexz(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("http://%s/indexz", TestConfig.BindAddress))
	if err != nil {
		t.Errorf("Error resolving bad URL - shoulvd return 404: %v", err)
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 for indexz")
	}
	if int(resp.ContentLength) != len(TestConfig.Content404) {
		t.Errorf("Expected length: %d, have length: %d", len(TestConfig.Content404), int(resp.ContentLength))
	}
}
