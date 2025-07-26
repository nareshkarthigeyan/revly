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
			{Role: "system", Content: `You're an expert Git user and code reviewer. Based on the following staged diff, write a clear and concise conventional commit message. You are given the diff of the staged changes, and you should focus on the changes made in this commit. DO NOT GIVE ANY EXPLANATION. Just the commit message.
			The whole commit message should be a single line, nothing more. Do not give multiple messages - put all of the messages into one. If the diff has so many changes, you should summarize the intent of the overall changes in a few words.
			You should follow the conventional commit format.
			Use the format: <type>(optional scope): <description>
			Types: feat, fix, refactor, docs, style, test, chore
			Here’s the diff:`},
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