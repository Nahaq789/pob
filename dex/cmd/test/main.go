package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gen "pob/dex/gen"
	"pob/pkg/interceptor/hmac"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	secret := os.Getenv("HMAC_SECRET")
	if secret == "" {
		log.Fatal("HMAC_SECRET is not set")
	}

	conn, err := grpc.NewClient("127.0.0.1:9091",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(hmac.HmacClientInterceptor(secret)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gen.NewDexServiceClient(conn)

	// GetPokemon
	start := time.Now()
	pokemon, err := client.GetPokemon(context.Background(), &gen.GetPokemonRequest{PokemonId: 1})
	if err != nil {
		log.Fatalf("GetPokemon: %v", err)
	}
	fmt.Printf("GetPokemon: %s (%v)\n", pokemon.Name, time.Since(start))

	// GetLearnableMoves
	start = time.Now()
	moves, err := client.GetLearnableMoves(context.Background(), &gen.GetLearnableMovesRequest{PokemonId: 1})
	if err != nil {
		log.Fatalf("GetLearnableMoves: %v", err)
	}
	fmt.Printf("GetLearnableMoves: %d moves (%v)\n", len(moves.Moves), time.Since(start))

	// GetMove
	start = time.Now()
	move, err := client.GetMove(context.Background(), &gen.GetMoveRequest{MoveId: 1})
	if err != nil {
		log.Fatalf("GetMove: %v", err)
	}
	fmt.Printf("GetMove: %s (%v)\n", move.Name, time.Since(start))

	// GetAbility
	start = time.Now()
	ability, err := client.GetAbility(context.Background(), &gen.GetAbilityRequest{AbilityId: 1})
	if err != nil {
		log.Fatalf("GetAbility: %v", err)
	}
	fmt.Printf("GetAbility: %s (%v)\n", ability.AbilityName, time.Since(start))

	// GetItem
	start = time.Now()
	item, err := client.GetItem(context.Background(), &gen.GetItemRequest{ItemId: 127})
	if err != nil {
		log.Fatalf("GetItem: %v", err)
	}
	fmt.Printf("GetItem: %s (%v)\n", item.Name, time.Since(start))
}
