package history

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Camilo-845/type-cli/internal/game"
	"github.com/Camilo-845/type-cli/internal/paths"
)

const MaxHistoryEntries = 100

func historyPath() (string, error) {
	dir, err := paths.AppDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "history.json"), nil
}

func Load() ([]game.Result, error) {
	path, err := historyPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []game.Result{}, nil
		}
		return nil, err
	}

	var results []game.Result
	if err := json.Unmarshal(data, &results); err != nil {
		return []game.Result{}, nil
	}

	return results, nil
}

func Save(results []game.Result) error {
	dir, err := paths.AppDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path, err := historyPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Append(result game.Result) error {
	results, err := Load()
	if err != nil {
		results = []game.Result{}
	}

	if len(results) >= MaxHistoryEntries {
		results = results[len(results)-MaxHistoryEntries+1:]
	}

	results = append(results, result)
	return Save(results)
}
