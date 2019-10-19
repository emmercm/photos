package hasher

import (
	"io"
	"os"
	"strconv"

	"github.com/OneOfOne/xxhash"
	"github.com/pkg/errors"
)

// SlowFile computes a hash of an entire file
func SlowFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file for slow hash")
	}
	defer f.Close()

	h := xxhash.New64()
	if _, err := io.Copy(h, f); err != nil {
		return "", errors.Wrap(err, "failed to slow hash file")
	}

	hash := h.Sum64()
	encoded := strconv.FormatUint(hash, 16)

	return encoded, nil
}
