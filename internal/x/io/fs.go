package iox

import (
	"encoding/json"
	"os"
)

func ReadJSON[T any](path string, o T) (T, error) {
	fs, err := os.ReadFile(path)
	if err != nil {
		return o, err
	}

	err = json.Unmarshal(fs, &o)
	if err != nil {
		return o, err
	}

	return o, err
}

func WriteJSON[T any](path string, o T) error {
	fs, err := json.Marshal(o)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, fs, 0o644)
	if err != nil {
		return err
	}

	return nil
}
