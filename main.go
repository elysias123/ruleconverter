package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	port  string
	debug bool
)

func main() {

	flag.StringVar(&port, "port", os.Getenv("RULCONVERTER_PORT"), "Port to run the server on")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	if port == "" {
		fmt.Printf("Port must be specified either via --port flag or RULCONVERTER_PORT environment variable")
		os.Exit(1)
	}

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	SetupRouter(r)

	fmt.Printf("ruleconverter is running...\n")
	fmt.Printf("Listen: http://localhost:%s\n", port)
	r.Run(fmt.Sprintf(":%s", port))
}
