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

package main

import (
	"os"

	bjorno "github.com/kris-nova/bjorn"

	"github.com/urfave/cli"

	"github.com/kris-nova/logger"
)

const (
	DefaultDefaultFile           = "index.html"
	DefaultInterpolateExtensions = ".html"
)

var (
	defaultFiles                 = cli.StringSlice{}
	interpolateExtensions        = cli.StringSlice{}
	statusPath404         string = ""
	statusPath500         string = ""
	statusPath5XX         string = ""
	useDefaultRootHandler bool   = false
)

func GetApp() *cli.App {
	app := &cli.App{
		Name: "DiaryApplication",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "verbose",
				Value:       4,
				Usage:       "verbosity level 0 (off) 4 (all)",
				Destination: &cfg.LogVerbosity,
			},
			&cli.StringFlag{
				Name:        "bind",
				Value:       ":80",
				Usage:       "default connection string",
				Destination: &cfg.BindAddress,
			},
			&cli.StringFlag{
				Name:        "dir",
				Value:       "/",
				Usage:       "default directory string",
				Destination: &cfg.ServeDirectory,
				EnvVar:      "BJORNODIR",
			},
			&cli.StringFlag{
				Name:        "notfound",
				Value:       "",
				Usage:       "default 404 not found file",
				Destination: &statusPath404,
				EnvVar:      "BJORNO404PATH",
			},
			&cli.StringFlag{
				Name:        "servererror",
				Value:       "",
				Usage:       "default 500 server error file",
				Destination: &statusPath500,
				EnvVar:      "BJORNO500PATH",
			},
			&cli.StringFlag{
				Name:        "servererrorall",
				Value:       "",
				Usage:       "default 5XX server error file",
				Destination: &statusPath5XX,
				EnvVar:      "BJORNO5XXPATH",
			},
			&cli.BoolFlag{
				Name:        "usedefault",
				Usage:       "use default filesystem handler",
				Destination: &cfg.UseDefaultRootHandler,
				EnvVar:      "BJORNOUSEDEFAULT",
			},
			&cli.StringSliceFlag{
				Name:  "default-files",
				Usage: "the default files to look for (index.html)",
				Value: &defaultFiles,
			},
			&cli.StringSliceFlag{
				Name:  "interpolate-extensions",
				Usage: "extensions to interpolate (.html)",
				Value: &defaultFiles,
			},
		},
		Action: func(c *cli.Context) error {
			// 404 handling
			if statusPath404 != "" {
				bytes, err := os.ReadFile(statusPath404)
				if err != nil {
					logger.Warning("Unable to load custom 404 path: %v", err)
					logger.Info("Using default 404 content.")
					cfg.Content404 = []byte(bjorno.StatusDefault404)
				} else {
					cfg.Content404 = bytes
				}
			}
			// 500 handling
			if statusPath500 != "" {
				bytes, err := os.ReadFile(statusPath500)
				if err != nil {
					logger.Warning("Unable to load custom 500 path: %v", err)
					logger.Info("Using default 500 content.")
					cfg.Content500 = []byte(bjorno.StatusDefault500)
				} else {
					cfg.Content500 = bytes
				}
			}
			// 5XX handling
			if statusPath5XX != "" {
				bytes, err := os.ReadFile(statusPath5XX)
				if err != nil {
					logger.Warning("Unable to load custom 5XX path: %v", err)
					logger.Info("Using default 5XX content.")
					cfg.Content5XX = []byte(bjorno.StatusDefault5XX)
				} else {
					cfg.Content5XX = bytes
				}
			}
			// Cast StringSlice to []string
			if len(interpolateExtensions) < 1 {
				interpolateExtensions.Set(DefaultInterpolateExtensions)
			}
			cfg.InterpolateExtensions = interpolateExtensions
			// Cast StringSlice to []string
			if len(defaultFiles) < 1 {
				defaultFiles.Set(DefaultDefaultFile)
			}
			cfg.DefaultIndexFiles = defaultFiles
			// ----------------------------------------
			return bjorno.Runtime(cfg, &bjorno.EmptyProgram{})
			// ----------------------------------------

		},
	}
	return app
}

var cfg = &bjorno.ServerConfig{}

func main() {
	app := GetApp()
	err := app.Run(os.Args)
	if err != nil {
		logger.Critical("Error running application: %v", err.Error())
	}
}
