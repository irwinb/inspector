package main

import (
	"github.com/irwinb/inspector/api"
	"github.com/irwinb/inspector/config"
	"github.com/irwinb/inspector/feeder"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting server.")

	log.Println("Initializing HTTP handlers.")

	http.Handle(config.ProxyEndpoint, http.StripPrefix(config.ProxyEndpoint,
		api.InspectorHandler(api.HandleProxy)))

	log.Println("Initializing feeder.")
	feeder.InitializeFeeder()

	log.Println("Starting feeder.")
	if err := http.ListenAndServe(config.ServerPort, nil); err != nil {
		log.Println("Staritng feeder failed: ", err)
	}
}
