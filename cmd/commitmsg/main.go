// CLI tool to get staged changes in a git repository and print them to the console
package main

import (
	"fmt"
	"os"

	"github.com/hazadus/gh-commitmsg/internal/git"
	"github.com/hazadus/gh-commitmsg/internal/llm"
)

func main() {
	stagedChanges, err := git.GetStagedChanges()
	if err != nil {
		fmt.Printf("Error retrieving staged changes: %v\n", err)
		os.Exit(1)
	}

	if stagedChanges == "" {
		fmt.Println("No staged changes in the repository.")
		return
	}

	llmClient, err := llm.NewClient()
	if err != nil {
		fmt.Printf("Error creating LLM client: %v\n", err)
		os.Exit(1)
	}

	commitMsg, err := llmClient.GenerateCommitMessage(stagedChanges)
	if err != nil {
		fmt.Printf("Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated commit message:")
	fmt.Println(commitMsg)
}

