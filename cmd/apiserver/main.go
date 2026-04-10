package main

import (
	"log"
	server "mifare/internal/app/apiserver"
	"mifare/internal/handler"
)

func main() {
	handlers := new(handler.Handler)

	srv := new(server.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
