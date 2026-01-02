package validator

import (
	"fmt"
	"os"
)

// IsValidMedia checks if the file is "healthy" enough to process
func IsValidMedia(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsPermission(err) {
			return false, fmt.Errorf("permission denied")
		}
		return false, err
	}

	// Rule 1: No empty files (prevents extraction crashes)
	if info.Size() == 0 {
		return false, fmt.Errorf("file is empty (0 bytes)")
	}

	// Rule 2: 100MB limit to protect MacBook memory
	const maxFileSize = 100 * 1024 * 1024
	if info.Size() > maxFileSize {
		return false, fmt.Errorf("file exceeds size limit (100MB)")
	}

	return true, nil
}
