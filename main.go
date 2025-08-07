package main

import (
	"incident-tracker/config"
	"incident-tracker/router"
	"incident-tracker/utils"
	"incident-tracker/workers"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	}
	c := config.LoadConfig()
	appContext := config.NewApplicationContext(c)
	server := router.NewServer(appContext)
	server.AddRoutes()
	actor := workers.CreateActors(appContext)
	actor.Start()
	defer actor.Stop()
	<-utils.WaitForTerminationHttpServer(server.Start())
}
