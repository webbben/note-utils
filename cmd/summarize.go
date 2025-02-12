package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webbben/note-utils/internal/llm"
	"github.com/webbben/note-utils/internal/util"
	"github.com/webbben/note-utils/pkg/summarize"
)

var maxHeader int

// summarizeCmd represents the summarize command
var summarizeCmd = &cobra.Command{
	Use:   "summarize",
	Short: "Generate a markdown summary of a given text",
	Long: `Generate a summary of a given text, in markdown format.
	For example:
	
	# generate summary for a file
	note-utils summarize --file todo.txt
	
	# generate a summary for content from stdin
	cat todo.txt | note-utils summarize`,
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

		opts := summarize.SummarizeOpts{}

		if fast {
			opts.Fast = true
		}
		if maxHeader > 1 {
			opts.MaxHeader = maxHeader
		}

		out, err := summarize.SummarizeNoteWithOpts(textContent, opts)
		if err != nil {
			return fmt.Errorf("error occurred while generating summary: %w", err)
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	summarizeCmd.Flags().BoolVar(&debug, "debug", false, "include debugging details in output")
	summarizeCmd.Flags().StringVarP(&file, "flag", "f", "", "specify a file path to summarize its contents.")
	summarizeCmd.Flags().BoolVar(&fast, "fast", false, "use a faster model for generating the summary. may result in lower quality, esp. for longer content.")
	summarizeCmd.Flags().IntVar(&maxHeader, "maxHeader", 0, "set the top header level for the markdown output; useful for when inserting summary content into a larger body of notes.")
	rootCmd.AddCommand(summarizeCmd)
}
