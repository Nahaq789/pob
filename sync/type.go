package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

const typeCsvPath = "../.data/types.csv"

type Type struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (t Type) GetId() int {
	return t.Id
}

type typeListResponse struct {
	Id    int `json:"id"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
}

type TypeRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewTypeRepository(client *ApiClient, db *DbClient) *TypeRepository {
	return &TypeRepository{c: client, db: db}
}

func (t *TypeRepository) fetch(ctx context.Context, id int) (*typeListResponse, error) {
	c := t.c.client
	url := t.c.baseURL.JoinPath(fmt.Sprintf("/type/%d/", id))
	method := http.MethodGet

	req, err := http.NewRequestWithContext(ctx, method, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var typeResp typeListResponse
	parseErr := json.NewDecoder(res.Body).Decode(&typeResp)
	if parseErr != nil {
		return nil, parseErr
	}

	return &typeResp, nil
}

func (t *TypeRepository) fetchAll(ctx context.Context) ([]*typeListResponse, error) {
	ids := make([]int, 18)
	for i := range ids {
		ids[i] = i + 1
	}
	return P(ctx, ids, t.fetch)
}

func (t *TypeRepository) toTypes(responses []*typeListResponse) []Type {
	var types []Type
	for _, res := range responses {
		for _, name := range res.Names {
			if name.Language.Name == "ja" {
				types = append(types, Type{Id: res.Id, Name: name.Name})
			}
		}
	}
	return types
}

func (t *TypeRepository) read() ([]Type, error) {
	return ReadCsv(string(typeCsvPath), func(s [][]string) ([]Type, error) {
		var types []Type
		for _, r := range s[1:] { // skip header
			id, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			types = append(types, Type{Id: id, Name: r[1]})
		}
		return types, nil
	})
}

func (t *TypeRepository) insert(ctx context.Context) error {
	types, err := t.read()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", typeCsvPath), slog.Int("count", len(types)))

	db := t.db.GetClient()
	for _, ty := range types {
		_, err := db.Exec(ctx,
			`INSERT INTO types (id, name) VALUES ($1, $2)`,
			ty.Id, ty.Name,
		)
		if err != nil {
			return err
		}
	}
	slog.InfoContext(ctx, "inserted types", slog.Int("count", len(types)))
	return nil
}

func (t *TypeRepository) write(ctx context.Context) error {
	responses, err := t.fetchAll(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched types", slog.Int("count", len(responses)))

	types := t.toTypes(responses)
	types = SortById(types)
	records := [][]string{
		{"id", "name"},
	}

	for _, ty := range types {
		records = append(records, []string{
			strconv.Itoa(ty.Id),
			ty.Name,
		})
	}

	if err := WriteCsv(typeCsvPath, records); err != nil {
		return err
	}

	slog.InfoContext(ctx, "wrote csv", slog.String("path", typeCsvPath))
	return nil
}

func (t *TypeRepository) ExecuteSync(ctx context.Context) error {
	if err := t.insert(ctx); err != nil {
		return err
	}
	return nil
}

func (t *TypeRepository) ExecuteCsv(ctx context.Context) error {
	if err := t.write(ctx); err != nil {
		return err
	}
	return nil
}
