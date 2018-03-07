package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sul-dlss-labs/taco"
	"github.com/sul-dlss-labs/taco/config"
	"github.com/sul-dlss-labs/taco/handlers"
)

func main() {
	rt, err := taco.NewRuntime(config.NewConfig())
	if err != nil {
		log.Fatalln(err)
	}

	server := createServer(rt)

	// serve API
	if err := server.Run(); err != nil {
		log.Fatalln(err)
	}
}

func createServer(rt *taco.Runtime) *gin.Engine {
	return handlers.BuildAPI(rt).Engine()
}
