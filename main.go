package main

import (
	"inspector/api"
	"inspector/config"
	"inspector/feeder"
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
