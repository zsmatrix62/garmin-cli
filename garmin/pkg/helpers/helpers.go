package helpers

import (
	"encoding/json"
	"fmt"
	"path"
)

func StateFileName(baseDir string, username string) string {
	return path.Join(baseDir, fmt.Sprintf("%s.json", username))
}

func JsonString(v interface{}) (string, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
