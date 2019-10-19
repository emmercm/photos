package hasher

import (
	"encoding/hex"

	"github.com/kalafut/imohash"
	"github.com/pkg/errors"
)

// FastFile computes a hash from a sample of a file
func FastFile(filename string) (string, error) {
	hash, err := imohash.SumFile(filename)
	if err != nil {
		return "", errors.Wrap(err, "failed to fast hash file")
	}

	encoded := hex.EncodeToString(hash[:])

	return encoded, nil
}
