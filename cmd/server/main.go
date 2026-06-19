package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"swe-zusatzuebung/internal/config"
	"swe-zusatzuebung/internal/database"
	"swe-zusatzuebung/internal/server"
)

func main() {
	cfg := config.Load()

	dbpool, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer dbpool.Close()

	log.Println("database connection is ready")

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("starting server on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, server.NewRouter()); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
