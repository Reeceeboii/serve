package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// global settings read in from cli globals
type settingsT struct {
	// what port will the server share over
	port string
	// what directory will be shared
	directory string
	// is verbose logging turned on
	verbose bool
	// will the directory be shared recursively (inc children)
	recursiveShare bool
}

// default settings
var settings = settingsT{"5000", ".", false, true}

// kick off the cli app with the program's command line arguments
func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// set up the server itself and serve the files at the directory
func serve(c *cli.Context) {
	root, err := getServerRoot(settings.directory)
	if err != nil {
		log.Fatal(err)
	}

	server := http.FileServer(http.Dir(root))
	http.Handle("/", server)
	if !settings.recursiveShare {

	}
	log.Println("Navigate to: 127.0.0.1:" + settings.port)
	err = http.ListenAndServe(":"+settings.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// get the absolute path of any directory given to the serve executable
func getServerRoot(pathArg string) (string, error) {
	return filepath.Abs(pathArg)
}

func nonRecursiveHandler(writer http.ResponseWriter, r *http.Request) string {
	return "go away"
}

// given a port, check if that port is open for our program to attach itself to
func isPortAvailable(port string) bool {
	if settings.verbose {
		log.Println("\t[-v]Checking if port " + port + " is available...")
	}

	ln, err := net.Listen("tcp", ":"+port)
	defer ln.Close()
	if err != nil {
		if settings.verbose {
			log.Println("\t[-v]Port " + port + " is not available\n")
		}
		return false
	}
	if settings.verbose {
		log.Println("\t[-v]Port " + port + " is available\n")
	}
	return true
}

// CLI app - parse command line args and do some setup
var app = cli.App{
	EnableBashCompletion: true,
	Name:                 "Serve",
	Usage:                "Turn any local directory into a static file server that be accessed from anywhere on you local network!",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "port",
			Value:    settings.port,
			Required: false,
			Aliases:  []string{"p"},
			Usage:    "The local port that the server is to listen on",
		},
		&cli.StringFlag{
			Name:     "directory",
			Value:    settings.directory,
			Required: true,
			Aliases:  []string{"d"},
			Usage:    "The directory to be served",
		},
		&cli.BoolFlag{
			Name:     "verbose",
			Value:    settings.verbose,
			Required: false,
			Aliases:  []string{"v"},
			Usage:    "Enable verbose logging",
		},
		&cli.BoolFlag{
			Name:     "recursive share",
			Value:    settings.recursiveShare,
			Required: false,
			Aliases:  []string{"r"},
			Usage:    "Enable recursive sharing (allows access to child directories of shared directory",
		},
	},
	// anon func fired at program launch
	Action: func(c *cli.Context) error {
		settings.port = c.String("port")
		settings.directory = c.String("directory")
		settings.verbose = c.Bool("verbose")
		settings.recursiveShare = c.Bool("recursive share")
		if !(isPortAvailable(settings.port)) {
			log.Println("Port " + settings.port + " is not able to listened on")
			os.Exit(1)
		}
		if settings.verbose {
			log.Println("\t[-v]Starting server...")
		}
		serve(c)
		return nil
	},
}
