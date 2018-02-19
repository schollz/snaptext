package main

import (
	"flag"
	"fmt"

	"github.com/schollz/snaptext/src"
)

var (
	doDebug bool
	port    string
)

func main() {
	flag.StringVar(&port, "port", "8002", "port to run server")
	flag.BoolVar(&doDebug, "debug", false, "enable debugging")
	flag.Parse()

	server.SetLogLevel("debug")

	err := src.Run(port)
	if err != nil {
		fmt.Printf("Error: '%s'", err.Error())
	}
}
