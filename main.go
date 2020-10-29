package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const (
	port      = 5000
	directory = "."
)

var app = cli.App{
	EnableBashCompletion: true,
	Name:                 "Serve",
	Usage:                "Turn any local directory into a static file server that be accessed from anywhere on you local network!",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "port",
			Value:    port,
			Required: false,
			Aliases:  []string{"p"},
			Usage:    "The local port that the server is to listen on",
		},
		&cli.StringFlag{
			Name:     "directory",
			Value:    directory,
			Required: true,
			Aliases:  []string{"d"},
			Usage:    "The directory to be served",
		},
	},
	Action: func(c *cli.Context) error {
		serve(c)
		return nil
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func serve(c *cli.Context) {
	root, err := getServerRoot(c.Args().Get(0))
	if err != nil {
		log.Fatal(err)
	}

	server := http.FileServer(http.Dir(root))
	http.Handle("/", server)
	log.Println("Server running on port " + c.String("port"))
	err = http.ListenAndServe(":"+c.String("port"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getServerRoot(pathArg string) (string, error) {
	return filepath.Abs(pathArg)
}
