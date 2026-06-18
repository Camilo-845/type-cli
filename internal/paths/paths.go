package paths

import (
	"os"
	"path/filepath"
)

func AppDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".typegame"), nil
}
