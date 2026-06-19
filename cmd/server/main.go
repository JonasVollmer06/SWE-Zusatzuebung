package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"swe-zusatzuebung/internal/config"
	"swe-zusatzuebung/internal/database"
	"swe-zusatzuebung/internal/fussballer"
	"swe-zusatzuebung/internal/server"
)

func main() {
	server.PrintBanner()

	cfg := config.Load()

	dbpool, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	sqlDB, err := dbpool.DB()
	if err != nil {
		log.Fatalf("database handle failed: %v", err)
	}
	defer sqlDB.Close()

	log.Println("database connection is ready")

	fussballerRepository := fussballer.NewRepository(dbpool)
	fussballerReadService := fussballer.NewReadService(fussballerRepository)
	fussballerWriteService := fussballer.NewWriteService(fussballerRepository)
	fussballerRouter := fussballer.NewRouter(fussballerReadService, fussballerWriteService)

	addr := fmt.Sprintf(":%s", cfg.Port)

	log.Printf("starting server on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, server.NewRouter(fussballerRouter)); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
