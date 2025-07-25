package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
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

	color.Yellow("=== BEGIN DIFF ===")
	color.White(string(diff))
	color.Yellow("=== END DIFF ===")
	color.Magenta("Diff length: %d bytes\n", len(diff))
	if len(diff) < 50 {
		return "", fmt.Errorf("diff too small or empty")
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " Thinking hard about your code..."
	s.Start()
	defer s.Stop()
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
				{Role: "system", Content: `You are Revly, a state-of-the-art AI code review assistant built by Naresh Karthigeyan. 
				You are acting as a highly experienced senior software engineer with deep expertise in modern software development practices. 
				Your task is to review Git code diffs with a focus on: Correctness, Performance, Readability, Maintainability, Security. 
				Provide clear, specific, and actionable feedback. Be friendly and constructive, but don’t hesitate to point out serious issues when necessary. 
				Speak as if you’re mentoring a peer, not criticizing a junior. 
				Use markdown formatting for code snippets and lists.
				Start the message by giving a kind greeting and a brief summary about the diff first -  not more than 150 words.
				Format each issue as:
				[SEVERITY] Line <line number>: <brief summary>
				Suggestion: <actionable recommendation>
				Explanation: <concise reasoning or tradeoff>
				Use one of the following severity levels: 
				[CRITICAL]: Functional bugs, security issues, or performance bottlenecks that must be fixed.
				[WARNING]: Bad practices, readability or maintainability concerns that should be addressed.
				[INFO]: Optional improvements, style suggestions, or minor clarity enhancements.
				Do not repeat or summarize the entire diff. Focus only on lines with actual issues or suggestions. 
				Don’t hallucinate context beyond what’s in the diff.
				If context is missing, point that out explicitly. You are not a general assistant. Only review the code. Do not explain what you are or engage in meta-discussion.End the review with a positive, concise summary if appropriate. Your goal is to help developers ship better code, faster, with confidence. Speak as if you’re mentoring a peer, not criticizing a junior. `,},
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