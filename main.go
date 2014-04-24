package main

import (
	"github.com/irwinb/inspector/api"
	"github.com/irwinb/inspector/feeder"
	"log"
)

func main() {
	log.Println("Starting server.")

	if err := api.InitAndListen(); err != nil {
		log.Println("Staritng feeder failed: ", err)
	}

	feeder.InitAndListen()
}
