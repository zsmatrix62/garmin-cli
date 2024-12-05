package helpers

import (
	"fmt"
	"path"
)

func StateFileName(baseDir string, username string) string {
	return path.Join(baseDir, fmt.Sprintf("%s.json", username))
}
