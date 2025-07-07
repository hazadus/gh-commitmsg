// CLI tool to get staged changes in a git repository and print them to the console
package main

import (
	"fmt"
	"os"

	"github.com/hazadus/gh-commitmsg/internal/git"
	"github.com/hazadus/gh-commitmsg/internal/llm"
	"github.com/spf13/cobra"
)

const extensionName = "standup"

var (
	flagLanguage string
)
var rootCmd = &cobra.Command{
	Use:   extensionName,
	Short: "Generate AI-powered commit messages",
	Long:  "A GitHub CLI extension that generates commit messages using GitHub Models and git repo data",
	RunE:  runCommitMsg,
}

func init() {
	rootCmd.Flags().StringVarP(&flagLanguage, "language", "l", "english", "Language to generate commit message in")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runCommitMsg(_ *cobra.Command, _ []string) error {
	stagedChanges, err := git.GetStagedChanges()
	if err != nil {
		return fmt.Errorf("failed to get staged changes: %w", err)
	}

	if stagedChanges == "" {
		fmt.Println("No staged changes in the repository.")
		return nil
	}

	llmClient, err := llm.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	commitMsg, err := llmClient.GenerateCommitMessage(stagedChanges, flagLanguage)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	fmt.Println("Generated commit message:")
	fmt.Println(commitMsg)
	return nil
}
