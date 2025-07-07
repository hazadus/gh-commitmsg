# gh-commitmsg

Based on [gh-standup](https://github.com/sgoedecke/gh-standup) by Sean Goedecke.

A GitHub CLI extension that generates AI-powered commit messages using staged git changes from current repository. It uses free [GitHub Models](https://docs.github.com/en/github-models) for inference, so you don't need to do any token setup - your existing GitHub CLI token will do fine!

## Installation

```bash
gh extension install hazadus/gh-commitmsg
gh commitmsg
```

### Organizations

To ensure the GitHub CLI can access your organization's data:

```bash
# Authenticate with GitHub CLI (if not already done)
gh auth login

# Authenticate with your organizations
gh auth refresh -h github.com -s read:org
```

### Prerequisites

- [GitHub CLI](https://cli.github.com/) installed and authenticated

## Usage

### Basic Usage

Generate a commit message for the staged changes:

```bash
gh commitmsg
```

### Advanced Options

```bash
# Use different language
gh commitmsg --language russian

# Use previous 3 commit messages as context
gh commitmsg --examples
```