package summarize

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/webbben/note-utils/internal/llm"
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
		out, err = llm.GenerateCompletion(noteContent, sysPrompt, "llama3.2:3b")
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
		out, err := llm.GenerateCompletion(noteContent, sysPrompt, "deepseek-r1")
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
