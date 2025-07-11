// CLI tool to generate AI-powered commit messages based on staged changes.
package main

import (
	"fmt"
	"os"

	"github.com/hazadus/gh-commitmsg/internal/git"
	"github.com/hazadus/gh-commitmsg/internal/llm"
	"github.com/spf13/cobra"
)

const extensionName = "commitmsg"

var (
	flagLanguage string
	flagExamples bool
	flagModel string
)
var rootCmd = &cobra.Command{
	Use:   extensionName,
	Short: "Generate AI-powered commit messages",
	Long:  "A GitHub CLI extension that generates commit messages using GitHub Models and git repo data",
	RunE:  runCommitMsg,
}

func init() {
	rootCmd.Flags().StringVarP(&flagLanguage, "language", "l", "english", "Language to generate commit message in")
	rootCmd.Flags().BoolVarP(&flagExamples, "examples", "e", false, "Add examples of commit messages to context")
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", "openai/gpt-4o", "GitHub Models model to use")
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

	// Add examples of previous commit messages to context if the flag is set
	latestCommitMessages := ""
	if flagExamples {
		latestCommitMessages, err = git.GetCommitMessages(3)
		if err != nil {
			return fmt.Errorf("failed to get latest commit messages: %w", err)
		}
		fmt.Println("  Adding examples of previous commit messages to context")
	}

	llmClient, err := llm.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	fmt.Println("  Language for commit message:", flagLanguage)
	commitMsg, err := llmClient.GenerateCommitMessage(stagedChanges, flagModel, flagLanguage, latestCommitMessages)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	fmt.Println("💬 Generated commit message:")
	fmt.Println(commitMsg)
	return nil
}
