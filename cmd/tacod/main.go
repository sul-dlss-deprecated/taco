package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/generated/restapi"
	"github.com/sul-dlss-labs/taco/handlers"
)

var portFlag = flag.Int("port", 8080, "Port to run this service on")

func main() {
	mode := os.Getenv("TACO_ENV")
	if mode == "" {
		mode = "development"
	}

	configFile := fmt.Sprintf("../../config/%s.yaml", mode)
	config.Init(configFile)

	rt, err := taco.NewRuntime(viper.GetViper())
	if err != nil {
		log.Fatalln(err)
	}

	server := createServer(rt)
	// serve API
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

func createServer(rt *taco.Runtime) *restapi.Server {
	server := restapi.NewServer(handlers.BuildAPI(rt))
	defer server.Shutdown()

	// parse flags
	flag.Parse()
	// set the port this service will be run on
	server.Port = *portFlag
	return server
}