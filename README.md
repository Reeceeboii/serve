# serve
## Turn any local directory into a static file server over your local network

Serve allows you to pick any directory on your machine, and expose its contents as a static file server. This means the 
contents of this directory will then be accessible by any other machine on your local network, or even anyone else on the 
internet if you decide to port forward the host system.

### Usage
`serve.exe [global options] command [command options] [arguments...]`

### Serve shell's current working directory:
* `> ./serve`

### Serve arbitrary directory
* `> ./serve -d /path/to/dir`

### Serve with recursive sharing off
* `> ./serve -d /path/to/dir -nr`

```
+-- file.txt
+-- /super_secret_directory
|   +-- begin-with-the-crazy-ideas.textile
|   +-- on-simplicity-in-technology.markdown

```
