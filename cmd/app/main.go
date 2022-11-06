package main

import (
	"gitflow/config"
	"gitflow/server"
	"log"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error")
	}
	a := server.NewApp(config)

	if err = a.Run(config.Port); err != nil {
	}
}
