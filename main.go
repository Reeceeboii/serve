package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// global settings read in from cli globals
type settingsType struct {
	// what port will the server share over
	port string
	// what directory will be shared
	directory string
	// is verbose logging turned on
	verbose bool
	// will the directory not be shared recursively (exc. children)
	nonRecursiveShare bool
}

// some info about the client
type clientInfoType struct {
	localOutboundIP net.IP
}

// default settings
var settings = settingsType{"5000", ".", false, false}

// set up client settings
var clientInfo = clientInfoType{nil}

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

	if settings.verbose {
		log.Println("\t[-v]Root directory is " + root)
	}

	mux := http.NewServeMux()
	server := http.FileServer(http.Dir(root))
	mux.Handle("/", logMiddleware(server))

	log.Println("\t   LOCAL | Navigate to: 127.0.0.1:" + settings.port)
	log.Println("\t NETWORK | Navigate to: " + clientInfo.localOutboundIP.String() + ":" + settings.port)

	err = http.ListenAndServe(":"+settings.port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

// get the absolute path of any directory given to the serve executable
func getServerRoot(pathArg string) (string, error) {
	return filepath.Abs(pathArg)
}

// logging
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if the request comes in for a directory and the user has enabled non recursive directory sharing,
		// we don't want the request to reach the file server
		if strings.HasSuffix(r.URL.Path, "/") && settings.nonRecursiveShare && len(r.URL.Path) > 1 {
			if settings.verbose {
				log.Println("\t[-v]Incoming request @ " + r.URL.Path + " DIR - BLOCKED w/ 404 response")
				http.NotFound(w, r)
				return
			}
		}
		if settings.verbose {
			log.Println("\t[-v]Incoming request @ " + r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
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
			log.Println("\t[-v]Port " + port + " is not available")
		}
		return false
	}
	if settings.verbose {
		log.Println("\t[-v]Port " + port + " is available")
	}
	return true
}

func getOutboundIPAddress() net.IP {
	if settings.verbose {
		log.Println("\t[-v]Querying local outbound IP address")
	}
	conn, err := net.Dial("udp", "1.2.3.4:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	outbound := conn.LocalAddr().(*net.UDPAddr)
	if settings.verbose {
		log.Println("\t[-v]Local outbound IP address is " + outbound.IP.String())
	}
	return outbound.IP
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
			Name:     "non-recursive",
			Value:    settings.nonRecursiveShare,
			Required: false,
			Aliases:  []string{"nr"},
			Usage:    "Disables recursive sharing (disallows access to child directories of shared root directory",
		},
	},
	// anon func fired at program launch
	Action: func(c *cli.Context) error {
		settings.port = c.String("port")
		settings.directory = c.String("directory")
		settings.verbose = c.Bool("verbose")
		settings.nonRecursiveShare = c.Bool("non-recursive")
		if !(isPortAvailable(settings.port)) {
			log.Println("Port " + settings.port + " is not able to listened on")
			os.Exit(1)
		}
		clientInfo.localOutboundIP = getOutboundIPAddress()
		if settings.verbose {
			log.Printf("\t[-v]Settings: %+v\n", settings)
			log.Printf("\t[-v]Client info: %+v\n", clientInfo)
			log.Println("\t[-v]Starting server\n...")
		}
		serve(c)
		return nil
	},
}
