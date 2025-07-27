// llm/env.go
package llm

import (
	// "fmt"
	// "os"
	"sync"

	"github.com/joho/godotenv"
)

var loadEnvOnce sync.Once

func loadEnv() {
	loadEnvOnce.Do(func() {
		// Load .env from current working directory or project root
		err := godotenv.Load()
		if err != nil {
			// fmt.Println("DEBUG: No .env file found or failed to load")
		} else {
			// fmt.Println("DEBUG: .env file loaded successfully")
		}
	})
}