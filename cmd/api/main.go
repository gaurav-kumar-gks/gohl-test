package main

import (
	"log"

	"gks.com/gohl-test/internal/config"
	"gks.com/gohl-test/internal/server"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to create server %v", err)
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server %v", err)
	}

}
