package main

import (
	"github.com/irwinb/inspector/api"
	"github.com/irwinb/inspector/feeder"
	"log"
)

func main() {
	log.Println("Starting server.")

	api.InitAndListen()

	feeder.InitAndListen()
}
