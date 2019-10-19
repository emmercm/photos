package monitor

import (
	"path/filepath"
	"strings"
)

func matchesGlob(path string) bool {
	pathExt := strings.ToLower(filepath.Ext(path))

	for _, ext := range [...]string{"jpg", "png"} {
		if pathExt == "."+ext {
			return true
		}
	}

	return false
}
