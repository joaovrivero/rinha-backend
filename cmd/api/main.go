package main

import (
	"log"
	"rinha-backend/internal/config"
	"rinha-backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	server := server.NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatal("Error running server:", err)
	}
}
