package json

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ParseJSON[T any](filePath string, data T) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("serialize to JSON: %w", err)
	}

	directory := filepath.Dir(filePath)

	if err := os.MkdirAll(directory, 0755); err != nil {
		return fmt.Errorf("cannot create directory '%s': %w", directory, err)
	}

	if err := os.WriteFile(filePath, jsonData, 0600); err != nil {
		return fmt.Errorf("cannot write to file '%s': %w", filePath, err)
	}

	return nil
}
