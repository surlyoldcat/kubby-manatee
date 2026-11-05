package main

import (
	"log"

	"kubmovie/api/internal/api"
	"kubmovie/api/internal/config"
	"kubmovie/api/internal/database"
	"kubmovie/api/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	router := api.NewRouter(repository.NewMovieRepository(db))
	log.Fatal(router.Run(":" + cfg.ServerPort))
}
