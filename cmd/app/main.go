package main

import (
	"gitflow/config"
	"gitflow/server"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	config, err := config.InitConfig()
	if err != nil {
		log.Fatalf("error")
	}
	logrus := logrus.Logger{}
	a := server.NewApp(config, logrus)

	if err = a.Run(); err != nil {
		log.Fatal(err)
	}
}
