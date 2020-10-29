package main

import (
	"fmt"

	"github.com/urfave/cli"
)

// App : the exposed command line interface
// uses github.com/urfave/cli
var App = cli.App{
	Name:  "Serve",
	Usage: "Turn any local directory into a static file server that be accessed from anywhere on you local network!",
	Action: func(c *cli.Context) error {
		fmt.Println("running!")
		return nil
	},
}
