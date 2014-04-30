package main

import (
	"github.com/irwinb/inspector/api"
	"github.com/irwinb/inspector/config"
	"github.com/irwinb/inspector/feeder"
	"log"
	"net/http"
)

func main() {

	log.Println("Initializing server.")
	api.Init()
	feeder.Init()

	log.Println("Listening...")

	if err := http.ListenAndServe(config.ServerPort, nil); err != nil {
		log.Println("Staritng server failed: ", err)
	}

}
