package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webbben/note-utils/internal/llm"
	"github.com/webbben/note-utils/internal/util"
	"github.com/webbben/note-utils/pkg/qa"
)

var chat bool

// qaCmd represents the qa command
var qaCmd = &cobra.Command{
	Use:   "qa",
	Short: "Ask a question about notes or other text content",
	Long: `Ask a question about notes or other text content. Optionally, start a Q&A chat session discussing the content of the text.
	Text content can be delivered by either stdin or by specifying a file path.
	
	Examples:
	
	# Ask a single question about text provided by file path
	note-utils qa --file notes.txt "What did I work on last week?"
	
	# Ask a question about text provided by stdin
	cat notes.txt | note-utils qa "What did I work on last week?"
	
	# Start a Q&A chat discussing the content of a text file
	note-utils qa --file notes.txt --chat`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			llm.DEBUG = true
		}

		textContent, err := util.ReadFileOrStdin(file)
		if err != nil {
			return err
		}
		if textContent == "" {
			return errors.New("error: no stdin or file flag detected")
		}

		opts := qa.AskNoteQuestionOpts{}
		if fast {
			opts.Fast = true
		}

		if chat {
			// do chat stuff
			return nil
		}

		if len(args) < 1 {
			return errors.New("error: no question provided")
		}

		out, err := qa.AskNoteQuestion(textContent, args[0], opts)
		if err != nil {
			return fmt.Errorf("error while generating answer: %w", err)
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	qaCmd.Flags().BoolVar(&debug, "debug", false, "include debugging details in output")
	qaCmd.Flags().BoolVar(&chat, "chat", false, "Start a chat session discussing the subject text content.")
	qaCmd.Flags().StringVarP(&file, "file", "f", "", "Specify a file of text content to serve as subject matter.")
	qaCmd.Flags().BoolVar(&fast, "fast", false, "use a faster model for generating the summary. may result in lower quality, esp. for longer content.")
	rootCmd.AddCommand(qaCmd)
}
