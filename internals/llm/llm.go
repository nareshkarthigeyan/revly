package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	// "fmt"
	"net/http"
	"os"

	"github.com/nareshkarthigeyan/revly/internals/config"
)
type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func GetLLMResponse(prompt string) (string, error) {
	// // DEBUG: Starting GetLLMResponse
	// fmt.Println("DEBUG: Starting GetLLMResponse")

	loadEnv() // Load environment variables
	
	cfg, err := config.GetConfig()
	if err != nil {
		// // DEBUG: Error getting config
		// fmt.Println("DEBUG: Error getting config:", err)
		return "", err
	}

	// // DEBUG: Config loaded
	// fmt.Printf("DEBUG: Config loaded: %+v\n", cfg)

	endpoint := cfg.LLM.Endpoint
	models := cfg.LLM.Models
	key := os.Getenv("LLM_API_KEY")
	if key == "" {
		// // DEBUG: Missing LLM_API_KEY
		// fmt.Println("DEBUG: Missing LLM_API_KEY")
		return "", errors.New("LLM_API_KEY not set")
	}

	// // DEBUG: Endpoint and models
	// fmt.Println("DEBUG: Endpoint:", endpoint)
	// fmt.Printf("DEBUG: Models: %v\n", models)

	for _, model := range models {
		// // DEBUG: Trying model
		// fmt.Println("DEBUG: Trying model:", model)

		reqBody := Request{
			Model: model,
			Messages: []Message{
				{
					Role: "system",
					Content: `You are an expert Git user. Your task is to output a concise, single-line, conventional commit message based on a provided Git diff.
Only return a commit message in the format: <type>(optional scope): <description>
NEVER provide explanation, suggestions, or additional lines. Do NOT add any text before or after the commit message. No multiple lines. No summaries.
Types: feat, fix, refactor, docs, style, test, chore
Here is the staged diff:`,
				},
				{Role: "user", Content: prompt},
			},
		}

		b, err := json.Marshal(reqBody)
		if err != nil {
			// // DEBUG: Error marshalling request
			// fmt.Printf("DEBUG: Error marshaling request for model %s: %v\n", model, err)
			continue
		}

		// DEBUG: Request body
		// fmt.Println("DEBUG: Request body:", string(b))

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
		if err != nil {
			// // DEBUG: Error creating request
			// fmt.Printf("DEBUG: Error creating request for model %s: %v\n", model, err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")


		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var out Response
		err = json.NewDecoder(resp.Body).Decode(&out)
		if err != nil {
			continue
		}

		if len(out.Choices) > 0 {
			return out.Choices[0].Message.Content, nil
		}
	}
	return "", errors.New("all LLM models failed to respond successfully")
}