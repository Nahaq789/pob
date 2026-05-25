package sync

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

const typeCsvPath = "../.data/types.csv"

type Type struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
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
	c *ApiClient
}

func NewTypeRepository(client *ApiClient) *TypeRepository {
	return &TypeRepository{c: client}
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

func (t *TypeRepository) WriteCsv(ctx context.Context) error {
	responses, err := t.fetchAll(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched types", slog.Int("count", len(responses)))

	types := t.toTypes(responses)
	records := [][]string{
		{"id", "name"},
	}

	for _, ty := range types {
		records = append(records, []string{
			strconv.Itoa(ty.Id),
			ty.Name,
		})
	}

	f, err := os.Create(typeCsvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(records)

	if err := w.Error(); err != nil {
		return err
	}
	slog.InfoContext(ctx, "wrote csv", slog.String("path", typeCsvPath))
	return nil
}
