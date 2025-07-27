package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
)

const cacheDir = ".revly/cache"

// Generate a cache key from input + model
func Key(diff []byte) string {
	h := sha256.New()
	h.Write(diff)           // write raw bytes
	h.Write([]byte(""))  // write model string
	return hex.EncodeToString(h.Sum(nil))
}

func Save(key string, content string) error {
	path := filepath.Join(cacheDir, key + ".txt")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func Load(key string) (string, bool) {
	path := filepath.Join(cacheDir, key + ".txt")
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}