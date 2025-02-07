package cleanup

import (
	"errors"
	"strings"

	"github.com/webbben/note-utils/internal/llm"
)

var systemPrompt = `
You are an assistant that formats a note to markdown, while possibly making some minor improvements to organize the content.
Preserve as much of the original wording as possible. If a note is already in markdown format, you may not need to make any changes.

Do NOT do summarize or re-word the note's content. Try to preserve the original text as much as possible.
Do NOT leave out any content or details from the original note.
Do NOT add any new content of your own (such as extra details or notes).

Do not output anything except the cleaned up note's content, since it will be piped directly into a function.

Don't forget... leave the wording almost identical to the original note. I'm serious!

Example 1)

Input:

Today I had a design meeting with Frank for the backend authentication system. 
He had some feedback he wanted to give, including that we should avoid 3rd party packages if possible, and also that the current proposal
might have some performance issues based on how we integrate with other services. He also added that an ERD would be nice.

We are planning to meet again next week, so I'll try to address some of these concerns for next time.

Output:

# Design Meeting for Backend Authentication System

Today I had a design meeting with Frank for the backend authentication system.
He had some feedback he wanted to give, including:

- we should avoid 3rd party packages if possible.
- the current proposal might have some performance issues based on how we integrate with other services.
- an ERD would be nice.

We are planning to meet again next week, so I'll try to address some of these concerns for next time.

Example 2)

Input:

I talked to Alex on the phone today about the Visa application process, and he had a lot to say. Apparently, you should NOT apply by mail.
The applications apparently go missing sometimes if you do... He also said that they requested a sample of his underwear, for some kind of identification process.
Weird if you ask me. Anyway, he's going to be flying in next week, so we will have to get everything wrapped up by then.

Follow-ups for next time: check visa status, call the immigration bureau about the underwear inspection status, and expedite shipping for his visa to arrive by next Thursday.

Output:

# Alex's Visa Application Status

I talked to Alex on the phone today about the Visa application process, and he had a lot to say. Apparently, **you should NOT apply by mail**.
The applications apparently go missing sometimes if you do... He also said that they requested a sample of his underwear, for some kind of identification process.
Weird if you ask me. Anyway, he's going to be flying in next week, so we will have to get everything wrapped up by then.

Follow-ups for next time:

- Check Visa status.
- Call the Immigration Bureau about the underwear inspection status.
- Expedite shipping for his Visa to arrive by next Thursday.
`

type CleanNoteOpts struct {
	Fast bool // Fast mode uses llama3.2:3b, which is faster but will yield lower quality (and possibly more paraphrased) results
}

func CleanNoteWithOpts(noteContent string, opts CleanNoteOpts) (string, error) {
	if opts.Fast {
		return llm.GenerateCompletion(noteContent, systemPrompt, "llama3.2:3b")
	}
	return CleanNoteCOT(noteContent)
}

func CleanNoteCOT(noteContent string) (string, error) {
	// sometimes deepseek seems to not close its thinking portion, so retry if that occurs.
	for range 3 {
		out, err := llm.GenerateCompletion(noteContent, systemPrompt, "deepseek-r1")
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
