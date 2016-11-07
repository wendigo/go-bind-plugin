package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

func fileChecksum(path string) (string, error) {
	hasher := sha256.New()

	file, err := os.Open(path)

	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
