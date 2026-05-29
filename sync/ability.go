package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"strconv"
)

type Ability struct {
	Id          int
	Name        string
	Description string
}

func (a Ability) GetId() int {
	return a.Id
}

type abilityApiResponse struct {
	Id    int `json:"id"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	FlavorTextEntries []struct {
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
		FlavorText string `json:"flavor_text"`
	} `json:"flavor_text_entries"`
}

type AbilityRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewAbilityRepository(client *ApiClient, db *DbClient) *AbilityRepository {
	return &AbilityRepository{c: client, db: db}
}

func (a *AbilityRepository) collectAbilityIds(ctx context.Context, gen int) ([]int, error) {
	r, ok := PokemonIdRangeByGen[gen]
	if !ok {
		return nil, fmt.Errorf("unsupported generation: %d", gen)
	}
	ids := make([]int, 0, (r.end-r.start)+1)
	for i := r.start; i <= r.end; i++ {
		ids = append(ids, i)
	}
	responses, err := P(ctx, ids, func(ctx context.Context, id int) (*pokemonApiResponse, error) {
		return fetchPokemon(ctx, a.c, id)
	})
	if err != nil {
		return nil, err
	}

	seen := make(map[int]struct{})
	var abilityIds []int
	for _, res := range responses {
		for _, ab := range res.Abilities {
			id := ExtractIdFromUrl(ab.Ability.Url)
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				abilityIds = append(abilityIds, id)
			}
		}
	}
	sort.Ints(abilityIds)
	return abilityIds, nil
}

func (a *AbilityRepository) fetch(ctx context.Context, id int) (*abilityApiResponse, error) {
	url := a.c.baseURL.JoinPath(fmt.Sprintf("/ability/%d/", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := a.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp abilityApiResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *AbilityRepository) toAbilities(responses []*abilityApiResponse) []Ability {
	var abilities []Ability
	for _, res := range responses {
		var jaName string
		for _, name := range res.Names {
			if name.Language.Name == "ja" {
				jaName = name.Name
				break
			}
		}

		var jaDesc string
		for _, entry := range res.FlavorTextEntries {
			if entry.Language.Name == "ja" {
				jaDesc = entry.FlavorText
			}
		}

		abilities = append(abilities, Ability{
			Id:          res.Id,
			Name:        jaName,
			Description: jaDesc,
		})
	}
	return abilities
}

func (a *AbilityRepository) read(gen int) ([]Ability, error) {
	return ReadCsv(a.buildSavePath(gen), func(s [][]string) ([]Ability, error) {
		var abilities []Ability
		for _, r := range s[1:] { // skip header
			id, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			if r[1] == "" {
				return nil, fmt.Errorf("name is empty at row: %v", r)
			}
			abilities = append(abilities, Ability{Id: id, Name: r[1], Description: r[2]})
		}
		return abilities, nil
	})
}

func (a *AbilityRepository) write(ctx context.Context, gen int) error {
	path := a.buildSavePath(gen)

	abilityIds, err := a.collectAbilityIds(ctx, gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "collected ability ids", slog.Int("count", len(abilityIds)))

	responses, err := P(ctx, abilityIds, a.fetch)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched abilities", slog.Int("count", len(responses)))

	abilities := a.toAbilities(responses)
	abilities = SortById(abilities)

	records := [][]string{{"id", "name", "description"}}
	for _, ab := range abilities {
		records = append(records, []string{
			strconv.Itoa(ab.Id),
			ab.Name,
			ab.Description,
		})
	}

	if err := WriteCsv(path, records); err != nil {
		return err
	}

	slog.InfoContext(ctx, "wrote csv", slog.String("path", path))
	return nil
}

func (a *AbilityRepository) insert(ctx context.Context, gen int) error {
	abilities, err := a.read(gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", a.buildSavePath(gen)), slog.Int("count", len(abilities)))

	db := a.db.GetClient()

	ids := make([]int, len(abilities))
	names := make([]string, len(abilities))
	descriptions := make([]string, len(abilities))
	for i, ab := range abilities {
		ids[i] = ab.Id
		names[i] = ab.Name
		descriptions[i] = ab.Description
	}

	_, err = db.Exec(ctx, `
		INSERT INTO abilities (id, name, description)
		SELECT * FROM UNNEST($1::int[], $2::text[], $3::text[])
	`, ids, names, descriptions)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "inserted abilities", slog.Int("count", len(abilities)))
	return nil
}

func (a *AbilityRepository) buildSavePath(gen int) string {
	return fmt.Sprintf("../.data/gen%d/abilities.csv", gen)
}

func (a *AbilityRepository) ExecuteSync(ctx context.Context, gen int) error {
	return a.insert(ctx, gen)
}

func (a *AbilityRepository) ExecuteCsv(ctx context.Context, gen int) error {
	return a.write(ctx, gen)
}
