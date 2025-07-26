package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
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
	key := os.Getenv("OPENROUTER_KEY")
	if key == "" {
		return "", errors.New("OPENROUTER_KEY not set")
	}

	reqBody := Request{
		Model: "mistralai/mistral-7b-instruct:free",
		Messages: []Message{
			{Role: "system", Content: `You are an expert Git user. Your task is to output a concise, single-line, conventional commit message based on a provided Git diff. 
			Only return a commit message in the format: <type>(optional scope): <description>
			NEVER provide explanation, suggestions, or additional lines. Do NOT add any text before or after the commit message. No multiple lines. No summaries.
			This is not a chat. Just return the one-line commit message. Nothing more.
			Types: feat, fix, refactor, docs, style, test, chore
			Here is the staged diff:`},
			{Role: "user", Content: prompt},
		},
	}

	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out Response
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("no response from LLM")
	}

	return out.Choices[0].Message.Content, nil
}