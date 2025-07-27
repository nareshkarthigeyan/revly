package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
)

const cacheDir = ".revly/cache"

func init() {
	os.MkdirAll(cacheDir, 0755)
}

func Key(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return hex.EncodeToString(hasher.Sum(nil))
}

func Save(key string, data []byte) error {
	cacheFile := filepath.Join(cacheDir, key)
	return ioutil.WriteFile(cacheFile, data, 0644)
}

func Load(key string) ([]byte, error) {
	cacheFile := filepath.Join(cacheDir, key)
	return ioutil.ReadFile(cacheFile)
}
