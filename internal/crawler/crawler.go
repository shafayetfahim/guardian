package crawler

import (
	"os"
	"path/filepath"
	"strings"
)

// Search finds files matching specific extensions (e.g., .jpg, .ARW)
func Search(root string, extensions []string) ([]string, error) {
	var matches []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		for _, ext := range extensions {
			if strings.HasSuffix(strings.ToLower(path), strings.ToLower(ext)) {
				matches = append(matches, path)
			}
		}
		return nil
	})

	return matches, err
}
