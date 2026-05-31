package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"pob/pkg/logger"
	pobsync "pob/sync"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

func main() {
	godotenv.Load()
	logger.InitLogger()

	root := &cobra.Command{Use: "sync"}
	root.AddCommand(newCsvCmd(), newSyncCmd())

	if err := root.Execute(); err != nil {
		log.Fatal(err)
	}
}

func newCsvCmd() *cobra.Command {
	var target string
	var gen int

	cmd := &cobra.Command{
		Use:   "csv",
		Short: "PokeAPIからCSVを生成する",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			noGenTargets := map[string]bool{"types": true, "items": true}
			if target == "" && gen == 0 {
				return fmt.Errorf("--gen は必須です")
			}
			if !noGenTargets[target] && target != "" && gen == 0 {
				return fmt.Errorf("--target %s には --gen の指定が必要です", target)
			}

			client, err := pobsync.NewApiClient(os.Getenv("POKEAPI_BASE_URL"))
			if err != nil {
				return err
			}

			switch target {
			case "types", "":
				return pobsync.NewTypeRepository(client, nil).ExecuteCsv(ctx)
			case "items":
				return pobsync.NewItemRepository(client, nil).ExecuteCsv(ctx)
			case "moves":
				return pobsync.NewMoveRepository(client, nil).ExecuteCsv(ctx, gen)
			case "abilities":
				return pobsync.NewAbilityRepository(client, nil).ExecuteCsv(ctx, gen)
			case "pokemon":
				return pobsync.NewPokemonRepository(client, nil).ExecuteCsv(ctx, gen)
			case "pokemon_abilities":
				return pobsync.NewPokemonAbilityRepository(client, nil).ExecuteCsv(ctx, gen)
			case "pokemon_moves":
				return pobsync.NewPokemonMoveRepository(client, nil).ExecuteCsv(ctx, gen)
			default:
				return fmt.Errorf("未実装のターゲット: %s", target)
			}
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "対象テーブル (types/items/moves/abilities/pokemon/pokemon_abilities/pokemon_moves)")
	cmd.Flags().IntVar(&gen, "gen", 0, "世代番号 (types/items以外は必須)")
	return cmd
}

func newSyncCmd() *cobra.Command {
	var target string
	var gen int

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "CSVからDBにINSERTする",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			if target == "" {
				return fmt.Errorf("--target は必須です (types/items/moves/abilities/pokemon/pokemon_abilities/pokemon_moves)")
			}
			noGenTargets := map[string]bool{"types": true, "items": true}
			if !noGenTargets[target] && gen == 0 {
				return fmt.Errorf("--target %s には --gen の指定が必要です", target)
			}

			db, err := pobsync.NewDbClient(ctx, os.Getenv("DB_DSN"))
			if err != nil {
				return err
			}
			defer db.Close()

			switch target {
			case "types":
				return pobsync.NewTypeRepository(nil, db).ExecuteSync(ctx)
			case "items":
				return pobsync.NewItemRepository(nil, db).ExecuteSync(ctx)
			case "moves":
				return pobsync.NewMoveRepository(nil, db).ExecuteSync(ctx, gen)
			case "abilities":
				return pobsync.NewAbilityRepository(nil, db).ExecuteSync(ctx, gen)
			case "pokemon":
				return pobsync.NewPokemonRepository(nil, db).ExecuteSync(ctx, gen)
			case "pokemon_abilities":
				return pobsync.NewPokemonAbilityRepository(nil, db).ExecuteSync(ctx, gen)
			case "pokemon_moves":
				return pobsync.NewPokemonMoveRepository(nil, db).ExecuteSync(ctx, gen)
			default:
				return fmt.Errorf("未実装のターゲット: %s", target)
			}
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "対象テーブル (types/items/moves/abilities/pokemon/pokemon_abilities/pokemon_moves)")
	cmd.Flags().IntVar(&gen, "gen", 0, "世代番号 (types/items以外は必須)")
	return cmd
}
