package sync

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
)

type PokemonAbility struct {
	PokemonId int
	AbilityId int
	Slot      int16
	IsHidden  bool
}

type PokemonAbilityRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewPokemonAbilityRepository(client *ApiClient, db *DbClient) *PokemonAbilityRepository {
	return &PokemonAbilityRepository{c: client, db: db}
}

func (p *PokemonAbilityRepository) fetch(ctx context.Context, id int) (*pokemonApiResponse, error) {
	return fetchPokemon(ctx, p.c, id)
}

func (p *PokemonAbilityRepository) fetchAll(ctx context.Context, gen int) ([]*pokemonApiResponse, error) {
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

func (p *PokemonAbilityRepository) toPokemonAbilities(responses []*pokemonApiResponse) []PokemonAbility {
	var abilities []PokemonAbility
	for _, res := range responses {
		for _, a := range res.Abilities {
			abilities = append(abilities, PokemonAbility{
				PokemonId: res.Id,
				AbilityId: ExtractIdFromUrl(a.Ability.Url),
				Slot:      int16(a.Slot),
				IsHidden:  a.IsHidden,
			})
		}
	}
	return abilities
}

func sortPokemonAbilities(abilities []PokemonAbility) {
	sort.Slice(abilities, func(i, j int) bool {
		if abilities[i].PokemonId != abilities[j].PokemonId {
			return abilities[i].PokemonId < abilities[j].PokemonId
		}
		return abilities[i].Slot < abilities[j].Slot
	})
}

func (p *PokemonAbilityRepository) read(gen int) ([]PokemonAbility, error) {
	return ReadCsv(p.buildSavePath(gen), func(s [][]string) ([]PokemonAbility, error) {
		var abilities []PokemonAbility
		for _, r := range s[1:] { // skip header
			pokemonId, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			abilityId, err := strconv.Atoi(r[1])
			if err != nil {
				return nil, err
			}
			slot, err := strconv.Atoi(r[2])
			if err != nil {
				return nil, err
			}
			isHidden, err := strconv.ParseBool(r[3])
			if err != nil {
				return nil, err
			}
			abilities = append(abilities, PokemonAbility{
				PokemonId: pokemonId,
				AbilityId: abilityId,
				Slot:      int16(slot),
				IsHidden:  isHidden,
			})
		}
		return abilities, nil
	})
}

func (p *PokemonAbilityRepository) write(ctx context.Context, gen int) error {
	path := p.buildSavePath(gen)
	responses, err := p.fetchAll(ctx, gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched pokemon for abilities", slog.Int("count", len(responses)))

	abilities := p.toPokemonAbilities(responses)
	sortPokemonAbilities(abilities)

	records := [][]string{
		{"pokemon_id", "ability_id", "slot", "is_hidden"},
	}
	for _, ab := range abilities {
		records = append(records, []string{
			strconv.Itoa(ab.PokemonId),
			strconv.Itoa(ab.AbilityId),
			strconv.Itoa(int(ab.Slot)),
			strconv.FormatBool(ab.IsHidden),
		})
	}

	if err := WriteCsv(path, records); err != nil {
		return err
	}

	slog.InfoContext(ctx, "wrote csv", slog.String("path", path))
	return nil
}

func (p *PokemonAbilityRepository) insert(ctx context.Context, gen int) error {
	abilities, err := p.read(gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", p.buildSavePath(gen)), slog.Int("count", len(abilities)))

	db := p.db.GetClient()

	pokemonIds := make([]int, len(abilities))
	abilityIds := make([]int, len(abilities))
	slots := make([]int, len(abilities))
	isHiddens := make([]bool, len(abilities))

	for i, ab := range abilities {
		pokemonIds[i] = ab.PokemonId
		abilityIds[i] = ab.AbilityId
		slots[i] = int(ab.Slot)
		isHiddens[i] = ab.IsHidden
	}

	_, err = db.Exec(ctx, `
		INSERT INTO pokemon_abilities (pokemon_id, ability_id, slot, is_hidden)
		SELECT * FROM UNNEST($1::int[], $2::int[], $3::int[], $4::bool[])
	`, pokemonIds, abilityIds, slots, isHiddens)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "inserted pokemon_abilities", slog.Int("count", len(abilities)))
	return nil
}

func (p *PokemonAbilityRepository) buildSavePath(gen int) string {
	return fmt.Sprintf("../.data/gen%d/pokemon_abilities.csv", gen)
}

func (p *PokemonAbilityRepository) ExecuteSync(ctx context.Context, gen int) error {
	return p.insert(ctx, gen)
}

func (p *PokemonAbilityRepository) ExecuteCsv(ctx context.Context, gen int) error {
	return p.write(ctx, gen)
}
