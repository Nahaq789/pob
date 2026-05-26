package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

var MoveIdRangeByGen = map[int]IdRange{
	1: {start: 1, end: 165},
}

type Move struct {
	Id          int
	Name        string
	TypeId      int
	Power       *int
	Accuracy    *int
	Pp          int
	DamageClass string
}

func (m Move) GetId() int {
	return m.Id
}

type moveResponse struct {
	Id    int `json:"id"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Type struct {
		Url string `json:"url"`
	} `json:"type"`
	DamageClass struct {
		Name string `json:"name"`
	} `json:"damage_class"`
	Power    *int `json:"power"`
	Accuracy *int `json:"accuracy"`
	Pp       int  `json:"pp"`
}

type MoveRepository struct {
	c  *ApiClient
	db *DbClient
}

func NewMoveRepository(client *ApiClient, db *DbClient) *MoveRepository {
	return &MoveRepository{c: client, db: db}
}

func (m *MoveRepository) fetch(ctx context.Context, id int) (*moveResponse, error) {
	c := m.c.client
	url := m.c.baseURL.JoinPath(fmt.Sprintf("/move/%d/", id))
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

	var moveResp moveResponse
	parseErr := json.NewDecoder(res.Body).Decode(&moveResp)
	if parseErr != nil {
		return nil, parseErr
	}
	return &moveResp, nil
}

func (m *MoveRepository) fetchAll(ctx context.Context, gen int) ([]*moveResponse, error) {
	r, ok := MoveIdRangeByGen[gen]
	if !ok {
		return nil, fmt.Errorf("unsupported generation: %d", gen)
	}
	ids := make([]int, 0, (r.end-r.start)+1)
	for i := r.start; i <= r.end; i++ {
		ids = append(ids, i)
	}
	return P(ctx, ids, m.fetch)
}

func (m *MoveRepository) toMoves(responses []*moveResponse) []Move {
	var moves []Move
	for _, res := range responses {
		var jaName string
		for _, name := range res.Names {
			if name.Language.Name == "ja" {
				jaName = name.Name
				break
			}
		}
		moves = append(moves, Move{
			Id:          res.Id,
			Name:        jaName,
			TypeId:      ExtractIdFromUrl(res.Type.Url),
			DamageClass: res.DamageClass.Name,
			Power:       res.Power,
			Accuracy:    res.Accuracy,
			Pp:          res.Pp,
		})
	}
	return moves
}

func (m *MoveRepository) read(gen int) ([]Move, error) {
	return ReadCsv(m.buildSavePath(gen), func(s [][]string) ([]Move, error) {
		var moves []Move
		for _, r := range s[1:] { // skip header
			id, err := strconv.Atoi(r[0])
			if err != nil {
				return nil, err
			}
			if r[1] == "" {
				return nil, fmt.Errorf("name is empty at row: %v", r)
			}
			typeId, err := strconv.Atoi(r[2])
			if err != nil {
				return nil, err
			}
			if r[3] == "" {
				return nil, fmt.Errorf("damage_class is empty at row: %v", r)
			}
			pp, err := strconv.Atoi(r[6])
			if err != nil {
				return nil, err
			}
			var power *int
			if r[4] != "" {
				v, err := strconv.Atoi(r[4])
				if err != nil {
					return nil, err
				}
				power = &v
			}
			var accuracy *int
			if r[5] != "" {
				v, err := strconv.Atoi(r[5])
				if err != nil {
					return nil, err
				}
				accuracy = &v
			}

			moves = append(moves, Move{
				Id:          id,
				Name:        r[1],
				TypeId:      typeId,
				DamageClass: r[3],
				Power:       power,
				Accuracy:    accuracy,
				Pp:          pp,
			})
		}
		return moves, nil
	})
}

func (m *MoveRepository) insert(ctx context.Context, gen int) error {
	moves, err := m.read(gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "read csv", slog.String("path", m.buildSavePath(gen)), slog.Int("count", len(moves)))

	db := m.db.GetClient()

	ids := make([]int, len(moves))
	names := make([]string, len(moves))
	typeIds := make([]int, len(moves))
	damageClasses := make([]string, len(moves))
	powers := make([]*int, len(moves))
	accuracies := make([]*int, len(moves))
	pps := make([]int, len(moves))

	for i, mv := range moves {
		ids[i] = mv.Id
		names[i] = mv.Name
		typeIds[i] = mv.TypeId
		damageClasses[i] = mv.DamageClass
		powers[i] = mv.Power
		accuracies[i] = mv.Accuracy
		pps[i] = mv.Pp
	}

	_, err = db.Exec(ctx, `
		INSERT INTO moves (id, name, type_id, damage_class, power, accuracy, pp)
		SELECT * FROM UNNEST($1::int[], $2::text[], $3::int[], $4::text[], $5::int[], $6::int[], $7::int[])
	`, ids, names, typeIds, damageClasses, powers, accuracies, pps)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "inserted moves", slog.Int("count", len(moves)))
	return nil
}

func (m *MoveRepository) write(ctx context.Context, gen int) error {
	moveCsvPath := m.buildSavePath(gen)
	responses, err := m.fetchAll(ctx, gen)
	if err != nil {
		return err
	}
	slog.InfoContext(ctx, "fetched moves", slog.Int("count", len(responses)))

	moves := m.toMoves(responses)
	moves = SortById(moves)
	records := [][]string{
		{"id", "name", "type_id", "damage_class", "power", "accuracy", "pp"},
	}

	for _, mv := range moves {
		power := ""
		if mv.Power != nil {
			power = strconv.Itoa(*mv.Power)
		}
		accuracy := ""
		if mv.Accuracy != nil {
			accuracy = strconv.Itoa(*mv.Accuracy)
		}
		records = append(records, []string{
			strconv.Itoa(mv.Id),
			mv.Name,
			strconv.Itoa(mv.TypeId),
			mv.DamageClass,
			power,
			accuracy,
			strconv.Itoa(mv.Pp),
		})
	}

	if err := WriteCsv(moveCsvPath, records); err != nil {
		return err
	}

	slog.InfoContext(ctx, "wrote csv", slog.String("path", moveCsvPath))
	return nil
}

func (m *MoveRepository) buildSavePath(gen int) string {
	return fmt.Sprintf("../.data/gen%d/moves.csv", gen)
}

func (m *MoveRepository) ExecuteSync(ctx context.Context, gen int) error {
	if err := m.insert(ctx, gen); err != nil {
		return err
	}
	return nil
}

func (m *MoveRepository) ExecuteCsv(ctx context.Context, gen int) error {
	if err := m.write(ctx, gen); err != nil {
		return err
	}
	return nil
}
