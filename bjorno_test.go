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
