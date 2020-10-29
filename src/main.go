package main

import (
	"log"
	"os"
)

func main() {
	err := App.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
