package main

import (
	"context"
	"log"
	"pob/pkg/logger"
	pobsync "pob/sync"
)

func main() {
	ctx := context.Background()
	logger.InitLogger()

	client, err := pobsync.NewApiClient("https://pokeapi.co/api/v2")
	if err != nil {
		log.Fatal(err)
	}

	typeRepo := pobsync.NewTypeRepository(client)
	if err := typeRepo.Write(ctx); err != nil {
		log.Fatal(err)
	}
}

