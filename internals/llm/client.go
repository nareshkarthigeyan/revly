package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func ReviewDiffWithLLM(diff string) (string, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("missing OPENROUTER_API_KEY env variable")
	}

	payload := map[string]interface{}{
		"model": "openrouter/claude-3-opus", // or any other model
		"messages": []Message{
			{Role: "system", Content: "You're a senior code reviewer. Be concise, direct, and highlight any issues in the diff."},
			{Role: "user", Content: "Please review the following git diff:\n\n" + diff},
		},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	choices := result["choices"].([]interface{})
	message := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)

	return message, nil
}