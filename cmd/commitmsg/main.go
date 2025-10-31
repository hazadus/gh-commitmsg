// CLI tool to generate AI-powered commit messages based on staged changes.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hazadus/gh-commitmsg/internal/git"
	"github.com/hazadus/gh-commitmsg/internal/llm"
	"github.com/spf13/cobra"
)

const (
	extensionName = "commitmsg"
	maxExamples   = 20
)

var (
	flagLanguage string
	flagExamples string
	flagModel    string
)
var rootCmd = &cobra.Command{
	Use:   extensionName,
	Short: "Generate AI-powered commit messages",
	Long:  "A GitHub CLI extension that generates commit messages using GitHub Models and git repo data",
	RunE:  runCommitMsg,
	Args:  handleArgs,
}

func init() {
	rootCmd.Flags().StringVarP(&flagLanguage, "language", "l", "english", "Language to generate commit message in")
	rootCmd.Flags().StringVarP(&flagExamples, "examples", "e", "", "Add N examples of commit messages to context (default 3 if flag is set without value, max 20)")
	rootCmd.Flags().Lookup("examples").NoOptDefVal = "3"
	rootCmd.Flags().StringVarP(&flagModel, "model", "m", "openai/gpt-4o", "GitHub Models model to use")
}

// handleArgs processes positional arguments to support --examples N syntax
func handleArgs(cmd *cobra.Command, args []string) error {
	// If --examples was set to its NoOptDefVal and there's a positional arg, use it
	if flagExamples == "3" && len(args) == 1 {
		// Check if examples flag was actually set
		if cmd.Flags().Changed("examples") {
			flagExamples = args[0]
			return nil
		}
	}
	// No positional arguments allowed normally
	if len(args) > 0 {
		return fmt.Errorf("unexpected argument: %s", args[0])
	}
	return nil
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

	// Parse and validate examples count
	examplesCount := 0
	if flagExamples != "" {
		examplesCount, err = strconv.Atoi(flagExamples)
		if err != nil {
			return fmt.Errorf("invalid examples count: %s", flagExamples)
		}
		if examplesCount < 1 || examplesCount > maxExamples {
			return fmt.Errorf("examples count must be between 1 and %d", maxExamples)
		}
	}

	// Add examples of previous commit messages to context if the flag is set
	latestCommitMessages := ""
	if examplesCount > 0 {
		latestCommitMessages, err = git.GetCommitMessages(examplesCount)
		if err != nil {
			return fmt.Errorf("failed to get latest commit messages: %w", err)
		}
		fmt.Printf("  Adding %d example(s) of previous commit messages to context\n", examplesCount)
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
