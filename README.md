# Serve: Turn any local directory into a static file server accessible over your local network

Serve allows you to pick any directory on your machine, and expose its contents as a static file server. This means the 
contents of this directory will then be accessible by any other machine on your local network, or even anyone else on the 
internet if you decide to port forward the host system.

## Usage
`serve [global options] command [command options] [arguments...]`

## Connections
The port used for connections depends on the provided flags. The default port is 5000
* From your own machine: `localhost:port` or `127.0.0.1:port`
* From elsewhere on your network: `192.168.0.25:port` where `192.168.0.25` is your machine's IP. Serve will attempt to 
tell you what this address is during launch.

## Help
* `> serve -h`

## Flags
```
--port value, -p value       The local port that the server is to listen on (default: "5000")
--directory value, -d value  The directory to be served (default: ".")
--verbose, -v                Enable verbose logging (default: false)
--non-recursive, --nr        Disables recursive sharing (disallows access to child directories of shared root directory (default: false)
--no-caching, --nc           Disables HTTP Cache-Control headers to allow browser caching of server responses (default: false)
--help, -h                   show help (default: false)
```

## Examples

### Serve shell's current working directory:
* `> ./serve`

### Serve arbitrary directory
* `> ./serve -d /path/to/dir`

### Serve with recursive sharing off
* `> ./serve -d /path/to/dir -nr`

```
+-- file.txt
+-- /super_secret_directory
|   +-- secret_plans.jpg
|   +-- evil_plan.pdf
+-- readme.md
```
*With non-recursive sharing on, only the top-level of the root share is exposed. For example, `file.txt` 
and `readme.md` will be shared, but any requests for `/super_secret_directory` or any of its children will be dropped
with a 404 response.*

### Serve with caching headers disabled
* `> ./serve -d /path/to/dir -nc`

*Disabling caching headers will mean that frequent requests from the same host will respond with a disk cache as opposed 
coming to the server as a new request*

### Using all flags
* `> serve -p 4500 -d ~/share -v -nr -nc`

*Serve, on local port 4500, the `~/share` directory, with verbose logging and non recursive sharing enabled but browser 
caching disabled.*