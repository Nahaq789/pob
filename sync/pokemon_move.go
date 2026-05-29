package sync

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
)

type PokemonMove struct {
	PokemonId int
	MoveId    int
}

type PokemonMoveRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewPokemonMoveRepository(client *ApiClient, db *DbClient) *PokemonMoveRepository {
	return &PokemonMoveRepository{c: client, db: db}
}

func (p *PokemonMoveRepository) fetch(ctx context.Context, id int) (*pokemonApiResponse, error) {
	return fetchPokemon(ctx, p.c, id)
}

func (p *PokemonMoveRepository) fetchAll(ctx context.Context, gen int) ([]*pokemonApiResponse, error) {
	r, ok := PokemonIdRangeByGen[gen]
	if !ok {
		return nil, fmt.Errorf("unsupported generation: %d", gen)
	}
	ids := make([]int, 0, (r.end-r.start)+1)
	for i := r.start; i <= r.end; i++ {
		ids = append(ids, i)
	}
	return P(ctx, ids, p.fetch)
}

func (p *PokemonMoveRepository) toPokemonMoves(responses []*pokemonApiResponse, gen int) []PokemonMove {
	r := MoveIdRangeByGen[gen]
	var moves []PokemonMove
	for _, res := range responses {
		for _, m := range res.Moves {
			moveId := ExtractIdFromUrl(m.Move.Url)
			if moveId >= r.start && moveId <= r.end {
				moves = append(moves, PokemonMove{PokemonId: res.Id, MoveId: moveId})
			}
		}
	}
	return moves
}

func sortPokemonMoves(moves []PokemonMove) {
	sort.Slice(moves, func(i, j int) bool {
		if moves[i].PokemonId != moves[j].PokemonId {
			return moves[i].PokemonId < moves[j].PokemonId
		}
		return moves[i].MoveId < moves[j].MoveId
	})
}

func (p *PokemonMoveRepository) read(gen int) ([]PokemonMove, error) {
	return ReadCsv(p.buildSavePath(gen), func(s [][]string) ([]PokemonMove, error) {
		var moves []PokemonMove
		for _, r := range s[1:] { // skip header
			pokemonId, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			moveId, err := strconv.Atoi(r[1])
			if err != nil {
				return nil, err
			}
			moves = append(moves, PokemonMove{PokemonId: pokemonId, MoveId: moveId})
		}
		return moves, nil
	})
}

func (p *PokemonMoveRepository) write(ctx context.Context, gen int) error {
	path := p.buildSavePath(gen)
	responses, err := p.fetchAll(ctx, gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched pokemon for moves", slog.Int("count", len(responses)))

	moves := p.toPokemonMoves(responses, gen)
	sortPokemonMoves(moves)

	records := [][]string{
		{"pokemon_id", "move_id"},
	}
	for _, mv := range moves {
		records = append(records, []string{
			strconv.Itoa(mv.PokemonId),
			strconv.Itoa(mv.MoveId),
		})
	}

	if err := WriteCsv(path, records); err != nil {
		return err
	}

	slog.InfoContext(ctx, "wrote csv", slog.String("path", path))
	return nil
}

func (p *PokemonMoveRepository) insert(ctx context.Context, gen int) error {
	moves, err := p.read(gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", p.buildSavePath(gen)), slog.Int("count", len(moves)))

	db := p.db.GetClient()

	pokemonIds := make([]int, len(moves))
	moveIds := make([]int, len(moves))

	for i, mv := range moves {
		pokemonIds[i] = mv.PokemonId
		moveIds[i] = mv.MoveId
	}

	_, err = db.Exec(ctx, `
		INSERT INTO pokemon_moves (pokemon_id, move_id)
		SELECT * FROM UNNEST($1::int[], $2::int[])
	`, pokemonIds, moveIds)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "inserted pokemon_moves", slog.Int("count", len(moves)))
	return nil
}

func (p *PokemonMoveRepository) buildSavePath(gen int) string {
	return fmt.Sprintf("../.data/gen%d/pokemon_moves.csv", gen)
}

func (p *PokemonMoveRepository) ExecuteSync(ctx context.Context, gen int) error {
	return p.insert(ctx, gen)
}

func (p *PokemonMoveRepository) ExecuteCsv(ctx context.Context, gen int) error {
	return p.write(ctx, gen)
}
