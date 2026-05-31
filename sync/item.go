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

const itemCsvPath = "../.data/items.csv"

var TargetItemCategoryIds = []int{3, 4, 5, 6, 7, 12, 13, 15, 17, 18, 19, 42, 44, 45, 46}

type Item struct {
	Id         int
	Name       string
	Category   string
	FlavorText string
}

func (i Item) GetId() int {
	return i.Id
}

type categoryItemsResponse struct {
	Items []struct {
		Url string `json:"url"`
	} `json:"items"`
}

type itemApiResponse struct {
	Id       int    `json:"id"`
	Category struct {
		Name string `json:"name"`
	} `json:"category"`
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
		Text string `json:"text"`
	} `json:"flavor_text_entries"`
}

type ItemRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewItemRepository(client *ApiClient, db *DbClient) *ItemRepository {
	return &ItemRepository{c: client, db: db}
}

func (r *ItemRepository) fetchCategory(ctx context.Context, id int) (*categoryItemsResponse, error) {
	url := r.c.baseURL.JoinPath(fmt.Sprintf("/item-category/%d/", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := r.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp categoryItemsResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *ItemRepository) collectItemIds(ctx context.Context) ([]int, error) {
	categoryResponses, err := P(ctx, TargetItemCategoryIds, r.fetchCategory)
	if err != nil {
		return nil, err
	}

	seen := make(map[int]struct{})
	var ids []int
	for _, cat := range categoryResponses {
		for _, item := range cat.Items {
			id := ExtractIdFromUrl(item.Url)
			if _, ok := seen[id]; !ok {
				seen[id] = struct{}{}
				ids = append(ids, id)
			}
		}
	}
	sort.Ints(ids)
	return ids, nil
}

func (r *ItemRepository) fetchItem(ctx context.Context, id int) (*itemApiResponse, error) {
	url := r.c.baseURL.JoinPath(fmt.Sprintf("/item/%d/", id))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), http.NoBody)
	if err != nil {
		return nil, err
	}
	res, err := r.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var resp itemApiResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (r *ItemRepository) toItems(responses []*itemApiResponse) []Item {
	var items []Item
	for _, res := range responses {
		var jaName string
		for _, n := range res.Names {
			if n.Language.Name == "ja" {
				jaName = n.Name
				break
			}
		}

		var jaFlavorText string
		for _, e := range res.FlavorTextEntries {
			if e.Language.Name == "ja" {
				jaFlavorText = e.Text
			}
		}

		items = append(items, Item{
			Id:         res.Id,
			Name:       jaName,
			Category:   res.Category.Name,
			FlavorText: jaFlavorText,
		})
	}
	return items
}

func (r *ItemRepository) write(ctx context.Context) error {
	itemIds, err := r.collectItemIds(ctx)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "collected item ids", slog.Int("count", len(itemIds)))

	responses, err := P(ctx, itemIds, r.fetchItem)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched items", slog.Int("count", len(responses)))

	items := r.toItems(responses)
	items = SortById(items)

	records := [][]string{{"id", "name", "category", "flavor_text"}}
	for _, item := range items {
		records = append(records, []string{
			strconv.Itoa(item.Id),
			item.Name,
			item.Category,
			Replacer(item.FlavorText),
		})
	}

	if err := WriteCsv(itemCsvPath, records); err != nil {
		return err
	}
	slog.InfoContext(ctx, "wrote csv", slog.String("path", itemCsvPath))
	return nil
}

func (r *ItemRepository) read() ([]Item, error) {
	return ReadCsv(itemCsvPath, func(s [][]string) ([]Item, error) {
		var items []Item
		for _, row := range s[1:] { // skip header
			id, err := strconv.Atoi(row[0])
			if err != nil {
				return nil, err
			}
			if row[1] == "" {
				return nil, fmt.Errorf("name is empty at row: %v", row)
			}
			items = append(items, Item{Id: id, Name: row[1], Category: row[2], FlavorText: row[3]})
		}
		return items, nil
	})
}

func (r *ItemRepository) insert(ctx context.Context) error {
	items, err := r.read()
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", itemCsvPath), slog.Int("count", len(items)))

	db := r.db.GetClient()

	ids := make([]int, len(items))
	names := make([]string, len(items))
	categories := make([]string, len(items))
	flavorTexts := make([]*string, len(items))
	for i, item := range items {
		ids[i] = item.Id
		names[i] = item.Name
		categories[i] = item.Category
		if item.FlavorText == "" {
			flavorTexts[i] = nil
		} else {
			ft := item.FlavorText
			flavorTexts[i] = &ft
		}
	}

	_, err = db.Exec(ctx, `
		INSERT INTO items (id, name, category, flavor_text)
		SELECT * FROM UNNEST($1::int[], $2::text[], $3::text[], $4::text[])
	`, ids, names, categories, flavorTexts)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "inserted items", slog.Int("count", len(items)))
	return nil
}

func (r *ItemRepository) ExecuteCsv(ctx context.Context) error {
	return r.write(ctx)
}

func (r *ItemRepository) ExecuteSync(ctx context.Context) error {
	return r.insert(ctx)
}
