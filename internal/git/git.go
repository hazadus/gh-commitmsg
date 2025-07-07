// Package git contains functions to interact with git repositories
package git

import (
	"fmt"
	"os/exec"
)

// GetStagedChanges executes the command git diff --staged and returns its output
func GetStagedChanges() (string, error) {
	// Check if we are in a git repository
	if !isGitRepository() {
		return "", fmt.Errorf("current directory is not a git repository")
	}

	// Execute the command git diff --staged
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing git diff --staged: %v", err)
	}

	return string(output), nil
}

// GetCommitMessages retrieves the last <count> commit messages from the git repository
func GetCommitMessages(count int) (string, error) {
	// Check if we are in a git repository
	if !isGitRepository() {
		return "", fmt.Errorf("current directory is not a git repository")
	}

	// Execute the command git log -n <count>
	cmd := exec.Command("git", "log", "-n", fmt.Sprintf("%d", count))
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error executing git log: %v", err)
	}

	return string(output), nil
}

// isGitRepository checks if the current directory is a git repository
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}
