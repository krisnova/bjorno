package bjorno

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/kris-nova/logger"
)

func TestRootExample(t *testing.T) {
	cfg := &ServerConfig{
		ServeDirectory: "bjorno.com",
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

func TestRequestPath(t *testing.T) {

	cases := map[string]string{
		"/":         "/",
		".":         "/",
		"//":        "/",
		"style.css": "/style.css",
		"/beeps":    "/beeps",
		"beeps":     "/beeps",
		"":          "/",
	}
	for input, expected := range cases {
		r := &http.Request{
			URL: &url.URL{
				Path: input,
			},
		}
		actual := RequestPath(r)
		if actual != expected {
			t.Errorf("actual (%s) != expected (%s)", actual, expected)
		}
	}
}

// TestFileDirectory will test FileDirectoryPath which will only
// succeed if a file is in fact found.
func TestFileDirectory(t *testing.T) {
	defaultFiles := []string{"index.html", "index.beeps", "index.boops"}
	httpDir := http.Dir("bjorno.com")
	happyCases := map[string]string{
		"/":            "index.html",
		"":             "index.html",
		"/beeps/boops": "boops",
		"/meeps":       "index.html",
	}
	// Should all be happy
	for input, expected := range happyCases {
		_, stat, err := FileDirectoryPath(defaultFiles, input, httpDir)
		if err != nil {
			t.Errorf("error FileDirectoryPath: %v", err)
		}
		actual := stat.Name()
		if actual != expected {
			t.Errorf("actual (%s) != expected (%s)", actual, expected)
		}
	}
	sadCases := []string{
		"beeps",
	}
	// Should all be sad
	for _, input := range sadCases {
		_, _, err := FileDirectoryPath(defaultFiles, input, httpDir)
		if err == nil {

			t.Errorf("expected err for input: %s", input)
		}
	}
}
