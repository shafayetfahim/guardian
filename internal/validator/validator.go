package validator

import (
	"fmt"
	"os"
)

func IsValidMedia(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsPermission(err) {
			return false, fmt.Errorf("permission denied")
		}
		return false, err
	}
	if info.Size() == 0 {
		return false, fmt.Errorf("file is empty (0 bytes)")
	}
	const maxFileSize = 100 * 1024 * 1024
	if info.Size() > maxFileSize {
		return false, fmt.Errorf("file exceeds size limit (100MB)")
	}
	return true, nil
}
