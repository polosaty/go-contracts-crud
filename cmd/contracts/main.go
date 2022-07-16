package main

import (
	"flag"
	"github.com/caarlos0/env/v6"

	"github.com/polosaty/go-contracts-crud/internal/app/config"
	"github.com/polosaty/go-contracts-crud/internal/app/server"
	"github.com/polosaty/go-contracts-crud/internal/app/storage"
	"log"
)

func main() {
	var cfg config.Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.RunAddress, "a", cfg.RunAddress, "server address")
	flag.StringVar(&cfg.DatabaseURI, "d", cfg.DatabaseURI, "database URI")
	flag.Parse()

	var db storage.Repository

	if db, err = storage.NewStoragePG(cfg.DatabaseURI); err != nil {
		log.Fatal(err)
	}
	log.Println("use postgres conn " + cfg.DatabaseURI + " as db")

	log.Fatal(server.Serve(cfg.RunAddress, db))
}
