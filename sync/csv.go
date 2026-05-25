package sync

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

func WriteCsv(path string, records [][]string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	w.WriteAll(records)
	return w.Error()
}
