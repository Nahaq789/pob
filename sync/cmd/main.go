package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"pob/pkg/logger"
	pobsync "pob/sync"
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

			if target != "types" && gen == 0 {
				return fmt.Errorf("--target %s には --gen の指定が必要です", target)
			}

			client, err := pobsync.NewApiClient(os.Getenv("POKEAPI_BASE_URL"))
			if err != nil {
				return err
			}

			switch target {
			case "types", "":
				return pobsync.NewTypeRepository(client, nil).ExecuteCsv(ctx)
			default:
				return fmt.Errorf("未実装のターゲット: %s", target)
			}
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "対象テーブル (types/moves/pokemon/pokemon_abilities/pokemon_moves)")
	cmd.Flags().IntVar(&gen, "gen", 0, "世代番号 (types以外は必須)")
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
				return fmt.Errorf("--target は必須です (types/moves/pokemon/pokemon_abilities/pokemon_moves)")
			}
			if target != "types" && gen == 0 {
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
			default:
				return fmt.Errorf("未実装のターゲット: %s", target)
			}
		},
	}

	cmd.Flags().StringVar(&target, "target", "", "対象テーブル (types/moves/pokemon/pokemon_abilities/pokemon_moves)")
	cmd.Flags().IntVar(&gen, "gen", 0, "世代番号 (types以外は必須)")
	return cmd
}
