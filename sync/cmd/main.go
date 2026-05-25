package main

import (
	"context"
	"log"
	pobsync "pob/sync"
)

func main() {
	ctx := context.Background()

	client, err := pobsync.NewApiClient("https://pokeapi.co/api/v2")
	if err != nil {
		log.Fatal(err)
	}

	typeRepo := pobsync.NewTypeRepository(client)
	if err := typeRepo.WriteCsv(ctx); err != nil {
		log.Fatal(err)
	}
}