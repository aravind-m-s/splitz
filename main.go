package main

import (
	"log"
	"splitz/config"
	"splitz/di"
)

func main() {
	config := config.InitConfig()

	server, err := di.InitServer(config)

	if err != nil {
		log.Fatal("Uanble to connect to db", err)
	} else {
		server.Start(config)
	}


}
