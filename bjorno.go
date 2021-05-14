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

	"github.com/kris-nova/logger"
)

const (
	StatusDefault404 string = `404 not found (bjorno)`
	StatusDefault500 string = `500 server error (bjorno)`
	StatusDefault5XX string = `5xx server error (bjorno)`
	EndpointRoot     string = "/"
)

// ServerConfig is the WebServer configuration component of bjorno
// This struct holds all of the WebServer bits (pun intended).
//
// We should only expose fields we would like a consumer of
// Bjorno to use.
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

	Endpoints []*Endpoint
}

// Runtime is exactly what you think it is. This is the runtime
// component of Bjorno and is most likely the place everything
// will go wrong.
//
// Here we have a stateless "Runtime" paradigm which means
// we NEVER want a .DoThing() workflow on Bjorno.
//
// You are either running Bjorno as a web server, or you
// are using Bjorno incorrectly.
//
// Here be dragons. Ye be warned.
//
//    cfg *ServerConfig    web server configuration for runtime serving.
//    V    RuntimeProgram  the top level runtime program to interpolate with.
//
// Note that V will have an extremely specific paradigm behind it.
// and by design we expect V to change at runtime. Bjorno *should
// be resilient enough to support a chaotic V.
//
// In other words, V should never break Bjorno. So have fun.
func Runtime(cfg *ServerConfig, V RuntimeProgram) error {
	// Dealing with legacy logger noise
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
		logger.BitwiseLevel = logger.LogEverything
	default:
		logger.BitwiseLevel = logger.LogEverything
	}

	// Debug
	logger.Info("ServeDirectory: %s", cfg.ServeDirectory)
	logger.Info("BindAddress: %s", cfg.BindAddress)
	logger.Info("Verbosity: %d", cfg.LogVerbosity)
	logger.Info("Log Level: %d", logger.BitwiseLevel)
	logger.Info("Default Index Files: %s", cfg.DefaultIndexFiles)

	// Endpoints
	for _, endpoint := range cfg.Endpoints {
		if endpoint.Pattern == "/" {
			return fmt.Errorf("Unable to use custom / root handler. Use the standard library if you want to do this. Go away.")
		}
		logger.Info("Registering endpoint: %s", endpoint.Pattern)
		http.Handle(endpoint.Pattern, endpoint.Handler)
	}

	logger.Info("Registering root endpoint: /")
	http.Handle("/", NewRootHandler(cfg, V))

	// Because we define custom handlers above we do not need to
	// pass in a "generic" handler here.
	return http.ListenAndServe(cfg.BindAddress, nil)
}

// RuntimeProgram is the program that Bjorno will run to interpolate your web pages.
// Have fun with your programs and be safe kiddos!
type RuntimeProgram interface {
	Values() interface{}
	Refresh()
	Lock()
	Unlock()
}

type Endpoint struct {
	Pattern string
	Handler http.Handler
}
