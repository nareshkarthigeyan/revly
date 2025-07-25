package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func ReviewDiffWithLLM(diff string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENROUTER_KEY not set in environment")
	}

	fmt.Println("Diff length:", len(diff))
	if len(diff) < 50 {
		return "", fmt.Errorf("Diff too small or empty")
	}

	models := []string{
		"qwen/qwen3-coder:free", // fallback
		"mistralai/mistral-7b-instruct:free",
		"openchat/openchat-3.5:free",
	}

	var lastErr error
	for _, model := range models {
		body := OpenRouterRequest{
			Model: model,
			Messages: []Message{
				{Role: "system", Content: "You are Revly, a state-of-the-art code review assistant developed by Naresh Karthigeyan. You are a highly experienced senior software engineer reviewing code for performance, correctness, readability, maintainability, and security. You provide precise, constructive, and actionable feedback on code diffs. Be friendly and compassionate. Focus solely on the code and how to improve it."},
				{Role: "user", Content: fmt.Sprintf("Please review this Git diff:\n\n%s", diff)},
			},
			Stream: false,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			lastErr = err
			continue
		}

		req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonBody))
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("OpenRouter-Referer", "https://github.com/nareshkarthigeyan/revly")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		bodyBytes, _ := io.ReadAll(resp.Body)
		// fmt.Println("📦 Raw LLM Response:", string(bodyBytes))

		var result OpenRouterResponse
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			lastErr = err
			continue
		}

		if len(result.Choices) == 0 {
			lastErr = fmt.Errorf("LLM returned no response for model %s", model)
			continue
		}

		return result.Choices[0].Message.Content, nil
	}

	return "", lastErr
}