# gh-commitmsg

Based on [gh-standup](https://github.com/sgoedecke/gh-standup) by Sean Goedecke.

A GitHub CLI extension that generates AI-powered [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/#summary) messages using staged git changes from current repository. It uses free [GitHub Models](https://docs.github.com/en/github-models) for inference, so you don't need to do any token setup - your existing GitHub CLI token will do fine!

If you find this utility useful, check my another tool – [gh-repomon](https://github.com/hazadus/gh-repomon).

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
# Stage changes in git repo
git add .

# Generate commit message
gh commitmsg
```

### Advanced Options

You can use and combine different options:

```bash
# Use different language
gh commitmsg --language russian

# Use previous commit messages as examples for LLM (default: 3, max: 20)
gh commitmsg --examples

# Or specify custom number of examples
gh commitmsg --examples 5

# Use a different AI model
gh commitmsg --model xai/grok-3-mini
```

### Output examples

LLM will generate something like:

```
feat: add CI/CD workflows, license, and update project details

- Introduced .github/workflows/ci.yml:
  - Added CI pipeline for testing, building, and artifact upload.
  - Configured matrix builds for multiple OS/architectures.
- Added .github/workflows/release.yml:
  - Automated release process triggered by version tags.
  - Builds binaries for Linux, macOS, and Windows.
```

For more examples, [see commit messages in this repo](https://github.com/hazadus/gh-commitmsg/commits/main/).
