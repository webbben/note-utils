package llm

import llama "github.com/webbben/ollama-wrapper"

func GenerateCompletion(input, systemPrompt, model string) (string, error) {
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
	return res, nil
}
