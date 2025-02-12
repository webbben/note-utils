package llm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	llama "github.com/webbben/ollama-wrapper"
)

var DEBUG = false

// GenerateCompletion
func GenerateCompletionCOT(input string, systemPrompt string) (string, error) {
	// sometimes deepseek seems to not close its thinking portion, so retry if that occurs.
	for range 3 {
		out, err := GenerateCompletion(input, systemPrompt, "deepseek-r1")
		if err != nil {
			return "", err
		}

		parts := strings.Split(out, "</think>")
		if len(parts) > 1 {
			return strings.TrimSpace(parts[1]), nil
		}
	}

	return "", errors.New("llm response did not have a </think> tag; attempted 3 times")
}

func GenerateCompletion(input, systemPrompt, model string) (string, error) {
	start := time.Now()

	if model == "" {
		model = "llama3.2:3b"
	}
	cmd, err := llama.StartServer()
	if err != nil {
		return "", err
	}
	defer llama.StopServer(cmd)

	llama.SetModel(model)

	client, err := llama.GetClient()
	if err != nil {
		return "", err
	}

	res, err := llama.GenerateCompletionWithOpts(client, systemPrompt, input, map[string]interface{}{
		"temperature": 0,
	})
	if err != nil {
		return "", err
	}

	if DEBUG {
		fmt.Println(time.Since(start))
	}
	return res, nil
}
