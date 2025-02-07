package summarize

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	llama "github.com/webbben/ollama-wrapper"
)

var systemPrompt = `
You are an assistant that summarizes notes. Given a body of text, summarize its contents and output it in a nice, cleanly formatted markdown document.

Do not output anything except the summarized note content, since it will be piped directly into a function.
`

type SummarizeOpts struct {
	MaxHeader int  // the maximum markdown header to use, e.g. 1 (#), 2 (##), 3 (###) etc. Default is 1.
	Fast      bool // if true, will use llama3.2:3b, which is faster. accuracy or quality may not be as good, especially for longer texts.
}

func SummarizeNoteWithOpts(noteContent string, opts SummarizeOpts) (string, error) {
	sysPrompt := systemPrompt

	var out string
	var err error
	if opts.Fast {
		out, err = SummarizeNote(noteContent, "llama3.2:3b", sysPrompt)
	} else {
		out, err = SummarizeNoteCOT(noteContent, sysPrompt)
	}
	if err != nil {
		return "", err
	}

	if opts.MaxHeader <= 1 {
		return out, nil
	}

	reader := bufio.NewReader(strings.NewReader(out))
	headerMod := ""
	var b strings.Builder
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break // end of output string
		}

		if headerMod == "" {
			// on the first line, check if top level header meets our requirements
			topHeaderSize := strings.Count(line, "#")
			if topHeaderSize >= opts.MaxHeader {
				return out, nil
			}
			headerMod = strings.Repeat("#", opts.MaxHeader-topHeaderSize)
		}

		// downgrade each header
		if headerMod == "" {
			log.Println("error: failed to detect header mod size")
			return out, nil
		}
		if line[0] == '#' {
			b.WriteString(fmt.Sprintf("%s%s", headerMod, line))
		} else {
			b.WriteString(line)
		}
	}
	return b.String(), nil
}

func SummarizeNoteCOT(noteContent string, sysPrompt string) (string, error) {
	// sometimes deepseek seems to not close its thinking portion, so retry if that occurs.
	for range 3 {
		out, err := SummarizeNote(noteContent, "deepseek-r1", sysPrompt)
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

func SummarizeNote(noteContent string, model string, sysPrompt string) (string, error) {
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

	res, err := llama.GenerateCompletionWithOpts(client, systemPrompt, noteContent, map[string]interface{}{
		"temperature": 0,
	})
	if err != nil {
		return "", err
	}
	return res, nil
}
