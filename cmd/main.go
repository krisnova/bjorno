package main

import (
	"os"

	"github.com/kris-nova/bjorno"

	"github.com/urfave/cli"

	"github.com/kris-nova/logger"
)

func GetApp() *cli.App {
	app := &cli.App{
		Name: "DiaryApplication",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "verbose",
				Value:       1,
				Usage:       "verbosity level 0/1 on/off",
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
			},
		},
		Action: func(c *cli.Context) error {
			return bjorno.RunServer(cfg)
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
