package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

var PokemonIdRangeByGen = map[int]IdRange{
	1: {start: 1, end: 151},
}

type Pokemon struct {
	Id            int
	Name          string
	Type1Id       int
	Type2Id       *int
	BaseHp        int
	BaseAttack    int
	BaseDefense   int
	BaseSpAttack  int
	BaseSpDefense int
	BaseSpeed     int
	WeightKg      float64
}

func (p Pokemon) GetId() int {
	return p.Id
}

// pokemonApiResponse is the shared response type for /pokemon/{id}.
// Used by PokemonRepository, PokemonAbilityRepository, and PokemonMoveRepository.
type pokemonApiResponse struct {
	Id     int `json:"id"`
	Weight int `json:"weight"`
	Types  []struct {
		Slot int `json:"slot"`
		Type struct {
			Url string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Abilities []struct {
		Ability struct {
			Url string `json:"url"`
		} `json:"ability"`
		Slot     int  `json:"slot"`
		IsHidden bool `json:"is_hidden"`
	} `json:"abilities"`
	Moves []struct {
		Move struct {
			Url string `json:"url"`
		} `json:"move"`
	} `json:"moves"`
}

type pokemonSpeciesApiResponse struct {
	Names []struct {
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
}

type pokemonFetchResult struct {
	pokemon *pokemonApiResponse
	species *pokemonSpeciesApiResponse
}

// fetchPokemon is a shared helper used by all three pokemon repositories.
func fetchPokemon(ctx context.Context, c *ApiClient, id int) (*pokemonApiResponse, error) {
	url := c.baseURL.JoinPath(fmt.Sprintf("/pokemon/%d/", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp pokemonApiResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

type PokemonRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewPokemonRepository(client *ApiClient, db *DbClient) *PokemonRepository {
	return &PokemonRepository{c: client, db: db}
}

func (p *PokemonRepository) fetch(ctx context.Context, id int) (*pokemonFetchResult, error) {
	pokemon, err := fetchPokemon(ctx, p.c, id)
	if err != nil {
		return nil, err
	}

	url := p.c.baseURL.JoinPath(fmt.Sprintf("/pokemon-species/%d/", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := p.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var species pokemonSpeciesApiResponse
	if err := json.NewDecoder(res.Body).Decode(&species); err != nil {
		return nil, err
	}

	return &pokemonFetchResult{pokemon: pokemon, species: &species}, nil
}

func (p *PokemonRepository) fetchAll(ctx context.Context, gen int) ([]*pokemonFetchResult, error) {
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

func (p *PokemonRepository) toPokemons(results []*pokemonFetchResult) []Pokemon {
	var pokemons []Pokemon
	for _, res := range results {
		var jaName string
		for _, name := range res.species.Names {
			if name.Language.Name == "ja" {
				jaName = name.Name
				break
			}
		}

		var type1Id int
		var type2Id *int
		for _, t := range res.pokemon.Types {
			id := ExtractIdFromUrl(t.Type.Url)
			if t.Slot == 1 {
				type1Id = id
			} else if t.Slot == 2 {
				type2Id = &id
			}
		}

		statsMap := make(map[string]int)
		for _, s := range res.pokemon.Stats {
			statsMap[s.Stat.Name] = s.BaseStat
		}

		pokemons = append(pokemons, Pokemon{
			Id:            res.pokemon.Id,
			Name:          jaName,
			Type1Id:       type1Id,
			Type2Id:       type2Id,
			BaseHp:        statsMap["hp"],
			BaseAttack:    statsMap["attack"],
			BaseDefense:   statsMap["defense"],
			BaseSpAttack:  statsMap["special-attack"],
			BaseSpDefense: statsMap["special-defense"],
			BaseSpeed:     statsMap["speed"],
			WeightKg:      float64(res.pokemon.Weight) / 10,
		})
	}
	return pokemons
}

func (p *PokemonRepository) read(gen int) ([]Pokemon, error) {
	return ReadCsv(p.buildSavePath(gen), func(s [][]string) ([]Pokemon, error) {
		var pokemons []Pokemon
		for _, r := range s[1:] { // skip header
			id, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			if r[1] == "" {
				return nil, fmt.Errorf("name is empty at row: %v", r)
			}
			type1Id, err := strconv.Atoi(r[2])
			if err != nil {
				return nil, err
			}
			var type2Id *int
			if r[3] != "" {
				v, err := strconv.Atoi(r[3])
				if err != nil {
					return nil, err
				}
				type2Id = &v
			}
			baseHp, err := strconv.Atoi(r[4])
			if err != nil {
				return nil, err
			}
			baseAttack, err := strconv.Atoi(r[5])
			if err != nil {
				return nil, err
			}
			baseDefense, err := strconv.Atoi(r[6])
			if err != nil {
				return nil, err
			}
			baseSpAttack, err := strconv.Atoi(r[7])
			if err != nil {
				return nil, err
			}
			baseSpDefense, err := strconv.Atoi(r[8])
			if err != nil {
				return nil, err
			}
			baseSpeed, err := strconv.Atoi(r[9])
			if err != nil {
				return nil, err
			}
			weightKg, err := strconv.ParseFloat(r[10], 64)
			if err != nil {
				return nil, err
			}
			pokemons = append(pokemons, Pokemon{
				Id:            id,
				Name:          r[1],
				Type1Id:       type1Id,
				Type2Id:       type2Id,
				BaseHp:        baseHp,
				BaseAttack:    baseAttack,
				BaseDefense:   baseDefense,
				BaseSpAttack:  baseSpAttack,
				BaseSpDefense: baseSpDefense,
				BaseSpeed:     baseSpeed,
				WeightKg:      weightKg,
			})
		}
		return pokemons, nil
	})
}

func (p *PokemonRepository) write(ctx context.Context, gen int) error {
	path := p.buildSavePath(gen)
	results, err := p.fetchAll(ctx, gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched pokemon", slog.Int("count", len(results)))

	pokemons := p.toPokemons(results)
	pokemons = SortById(pokemons)

	records := [][]string{
		{"id", "name", "type1_id", "type2_id", "base_hp", "base_attack", "base_defense", "base_sp_attack", "base_sp_defense", "base_speed", "weight_kg"},
	}
	for _, pk := range pokemons {
		type2Id := ""
		if pk.Type2Id != nil {
			type2Id = strconv.Itoa(*pk.Type2Id)
		}
		records = append(records, []string{
			strconv.Itoa(pk.Id),
			pk.Name,
			strconv.Itoa(pk.Type1Id),
			type2Id,
			strconv.Itoa(pk.BaseHp),
			strconv.Itoa(pk.BaseAttack),
			strconv.Itoa(pk.BaseDefense),
			strconv.Itoa(pk.BaseSpAttack),
			strconv.Itoa(pk.BaseSpDefense),
			strconv.Itoa(pk.BaseSpeed),
			strconv.FormatFloat(pk.WeightKg, 'f', 1, 64),
		})
	}

	if err := WriteCsv(path, records); err != nil {
		return err
	}

	slog.InfoContext(ctx, "wrote csv", slog.String("path", path))
	return nil
}

func (p *PokemonRepository) insert(ctx context.Context, gen int) error {
	pokemons, err := p.read(gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", p.buildSavePath(gen)), slog.Int("count", len(pokemons)))

	db := p.db.GetClient()

	ids := make([]int, len(pokemons))
	names := make([]string, len(pokemons))
	type1Ids := make([]int, len(pokemons))
	type2Ids := make([]*int, len(pokemons))
	baseHps := make([]int, len(pokemons))
	baseAttacks := make([]int, len(pokemons))
	baseDefenses := make([]int, len(pokemons))
	baseSpAttacks := make([]int, len(pokemons))
	baseSpDefenses := make([]int, len(pokemons))
	baseSpeeds := make([]int, len(pokemons))
	weightKgs := make([]float64, len(pokemons))

	for i, pk := range pokemons {
		ids[i] = pk.Id
		names[i] = pk.Name
		type1Ids[i] = pk.Type1Id
		type2Ids[i] = pk.Type2Id
		baseHps[i] = pk.BaseHp
		baseAttacks[i] = pk.BaseAttack
		baseDefenses[i] = pk.BaseDefense
		baseSpAttacks[i] = pk.BaseSpAttack
		baseSpDefenses[i] = pk.BaseSpDefense
		baseSpeeds[i] = pk.BaseSpeed
		weightKgs[i] = pk.WeightKg
	}

	_, err = db.Exec(ctx, `
		INSERT INTO pokemon (id, name, type1_id, type2_id, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, weight_kg)
		SELECT * FROM UNNEST($1::int[], $2::text[], $3::int[], $4::int[], $5::int[], $6::int[], $7::int[], $8::int[], $9::int[], $10::int[], $11::numeric[])
	`, ids, names, type1Ids, type2Ids, baseHps, baseAttacks, baseDefenses, baseSpAttacks, baseSpDefenses, baseSpeeds, weightKgs)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "inserted pokemon", slog.Int("count", len(pokemons)))
	return nil
}

func (p *PokemonRepository) buildSavePath(gen int) string {
	return fmt.Sprintf("../.data/gen%d/pokemon.csv", gen)
}

func (p *PokemonRepository) ExecuteSync(ctx context.Context, gen int) error {
	return p.insert(ctx, gen)
}

func (p *PokemonRepository) ExecuteCsv(ctx context.Context, gen int) error {
	return p.write(ctx, gen)
}
