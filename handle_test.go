// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
