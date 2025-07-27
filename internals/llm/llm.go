package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	loadEnv() // Load environment variables

	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	endpoint := cfg.LLM.Endpoint
	models := cfg.LLM.Models
	key := os.Getenv("LLM_API_KEY")
	if key == "" {
		return "", errors.New("LLM_API_KEY not set")
	}

	for _, model := range models {
		reqBody := Request{
			Model: model,
			Messages: []Message{
				{
					Role: "system",
					Content: `You are an expert Git user. Your task is to output a concise, single-line, conventional commit message based on a provided Git diff.\nOnly return a commit message in the format: <type>(optional scope): <description>\nNEVER provide explanation, suggestions, or additional lines. Do NOT add any text before or after the commit message. No multiple lines. No summaries.\nTypes: feat, fix, refactor, docs, style, test, chore\nHere is the staged diff:`, 
				},
				{Role: "user", Content: prompt},
			},
		}

		b, err := json.Marshal(reqBody)
		if err != nil {
			continue
		}

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
		if err != nil {
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

func GetPairProgrammingComment(diff string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}

	endpoint := cfg.LLM.Endpoint
	models := cfg.LLM.Models
	key := os.Getenv("LLM_API_KEY")
	if key == "" {
		return "", errors.New("LLM_API_KEY not set")
	}

	for _, model := range models {
		reqBody := Request{
			Model: model,
			Messages: []Message{
				{
					Role:    "system",
					Content: `You are an expert software engineer acting as a **pair programming assistant**. Your job is to review small code changes in real-time as a developer is typing.

You will receive a code diff from a Git working directory. It will look like:

diff
diff --git a/main.go b/main.go
index 1a2b3c..4d5e6f 100644
--- a/main.go
+++ b/main.go
@@ func myFunction() {
-   result := doSomething(x, y)
+   result := doSomething(x, y)
+   log.Printf("Result: %v", result)
}
Your job is to return a single, concise, high-signal suggestion or observation — max 20 words. This is not a full code review, but a lightweight, contextual nudge like a coding partner might give.

Focus on:
	•	readability
	•	naming
	•	small bug risks
	•	redundant logic
	•	performance
	•	unused code
	•	missing edge cases
	•	security issues
	•	any other small improvements
Avoid:
	•	large architectural changes
	•	major refactors
	•	overly complex suggestions
	•	anything that requires deep context beyond the diff
	•	anything that would require a full code review

					Format your response as a single line comment, like this:
					Have a brief greeting, or a follow up question if appropriate.
					Keep it short, actionable, and relevant to the diff provided. No explanations, just the comment itself.
					Don’t suggest changes that are already present in the diff.
`,
				},
				{Role: "user", Content: fmt.Sprintf("Code:\n%s", diff)},
			},
		}

		b, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Printf("Error marshaling request for model %s: %v\n", model, err)
			continue
		}

		req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
		if err != nil {
			fmt.Printf("Error creating request for model %s: %v\n", model, err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error calling LLM model %s: %v\n", model, err)
			continue
		}
		defer resp.Body.Close()

		var out Response
		err = json.NewDecoder(resp.Body).Decode(&out)
		if err != nil {
			fmt.Printf("Error decoding response from model %s: %v\n", model, err)
			continue
		}

		if len(out.Choices) > 0 {
			return out.Choices[0].Message.Content, nil
		}
		fmt.Printf("Model %s returned no choices.\n", model)
	}

	return "", errors.New("all LLM models failed to respond successfully")
}
