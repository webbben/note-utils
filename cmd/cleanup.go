package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/webbben/note-utils/internal/llm"
	"github.com/webbben/note-utils/internal/util"
	"github.com/webbben/note-utils/pkg/cleanup"
)

// cleanupCmd represents the cleanup command
var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up notes, organizing it for better readability and formatting it as markdown.",
	Long: `Clean up notes, organizing it for betterr readability and formatting it as markdown.
	Only accepts text content from stdin.
	
	# clean up notes from stdin, and save to a file
	cat notes.txt | note-utils cleanup > cleaned-notes.md`,
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

		opts := cleanup.CleanNoteOpts{}

		if fast {
			opts.Fast = true
		}

		out, err := cleanup.CleanNoteWithOpts(textContent, opts)
		if err != nil {
			return fmt.Errorf("error occurred while generating summary: %w", err)
		}
		fmt.Println(out)
		return nil
	},
}

func init() {
	cleanupCmd.Flags().BoolVar(&fast, "fast", false, "use a faster model for generating the summary. may result in lower quality, esp. for longer content.")
	cleanupCmd.Flags().BoolVar(&debug, "debug", false, "include debugging details in output")
	rootCmd.AddCommand(cleanupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
