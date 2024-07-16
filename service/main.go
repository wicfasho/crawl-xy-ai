package main

import (
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/wicfasho/crawl-xy/api/v1"
	"github.com/wicfasho/crawl-xy/database"
	"github.com/wicfasho/crawl-xy/scrapper"
	"github.com/wicfasho/crawl-xy/sqlc"
)

var wg sync.WaitGroup

func main() {
	// Initialize the database
	if err := database.Init(); err != nil {
		log.Fatal().Err(err).Msg("error initializing the database")
	}

	// Create DB Instance
	dbStore := sqlc.New(database.GetDB())

	// Start Scrapper Service
	wg.Add(1)
	go func() {
		scrapper.Start()
		wg.Done()
	}()

	// Spin up API Server
	runAPISever(dbStore)
}

func runAPISever(dbStore *sqlc.Queries) error {
	api, err := api.NewServer(dbStore)
	if err != nil {
		return err
	}

	api.Start(&wg)

	return nil
}
