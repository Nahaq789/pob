package main

import (
	"context"
	"log"
	"os"
	"pob/pkg/logger"
	pobsync "pob/sync"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load()
	logger.InitLogger()

	client, err := pobsync.NewApiClient(os.Getenv("POKEAPI_BASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := pobsync.NewDbClient(ctx, os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	typeRepo := pobsync.NewTypeRepository(client, db)
	if err := typeRepo.ExecuteCsv(ctx); err != nil {
		log.Fatal(err)
	}
	if err := typeRepo.ExecuteSync(ctx); err != nil {
		log.Fatal(err)
	}
}
