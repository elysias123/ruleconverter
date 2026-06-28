package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	port       string
	listenAddr string
	debug      bool
)

func main() {

	defaultPort := os.Getenv("RULCONVERTER_PORT")
	if defaultPort == "" {
		defaultPort = "30000"
	}
	flag.StringVar(&port, "port", defaultPort, "Port to run the server on")
	flag.StringVar(&listenAddr, "listen", "0.0.0.0", "Address to listen on")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.SetPrefix("[ruleconverter] ")
	r := gin.Default()
	SetupRouter(r)

	log.Printf("ruleconverter is running...\n")
	log.Printf("Listen: http://%s:%s\n", listenAddr, port)
	r.Run(fmt.Sprintf("%s:%s", listenAddr, port))
}
