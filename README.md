# Claude Commit Action Test Repository

This repository is created for testing the Claude PR Assistant GitHub Action.

## 📁 .github Directory Structure

```
.github/
├── PULL_REQUEST_TEMPLATE.md
└── workflows/
    ├── claude_pr_description.yml
    ├── claude_pr_review.yml
    └── root_workflows.yml
```

## 🔧 GitHub Actions Workflows

### 1. claude_pr_description.yml - PR Description Formatter

A reusable workflow that automatically formats PR descriptions using Claude AI.

**Key Features:**

- Analyzes PR changes to generate template-compliant descriptions
- Automatically fills Pull Request templates
- Uses `anthropics/claude-code-action@beta` action

**Required Secrets:**

- `ANTHROPIC_API_KEY`: Anthropic API key for Claude API access

**Allowed Tools:**

- `Bash(gh:*)`: GitHub CLI commands
- `Read`: File reading
- `View`: File viewing

### 2. claude_pr_review.yml - PR Code Reviewer

A reusable workflow that provides automated code review using Claude AI.

**Key Features:**

- Reviews PR code changes and provides inline feedback
- Uses GitHub's review system for structured comments
- Focuses on code quality, security, performance, and maintainability
- Submits non-blocking reviews (COMMENT type)

**Required Secrets:**

- `ANTHROPIC_API_KEY`: Anthropic API key for Claude API access

**Allowed Tools:**

- `mcp__github__create_pending_pull_request_review`: Start a review
- `mcp__github__add_comment_to_pending_review`: Add inline comments
- `mcp__github__submit_pending_pull_request_review`: Submit the review
- `mcp__github__get_pull_request_diff`: Get diff information

### 3. root_workflows.yml - Main Workflow Orchestrator

The main workflow that calls both Claude PR Assistant workflows.

**Trigger Conditions:**

- When a PR is opened (`pull_request.opened`) with changes to `.github/**` files
- When `@claude_pr` is mentioned in issue comments

**Jobs:**

- `call-claude-pr-description`: Calls the PR description formatter
- `call-claude-default`: Calls the PR code reviewer

**Permissions:**

- `contents: read` - Read repository contents
- `pull-requests: write` - Modify PRs
- `issues: write` - Write issues
- `id-token: write` - Write ID tokens

## 🚀 Usage

1. **On PR Creation**: When creating a new PR that modifies `.github/**` files:

   - Claude automatically formats the PR description according to the template
   - Claude provides code review feedback with inline comments

2. **Manual Trigger**: Mention `@claude_pr` in PR comments to manually trigger both workflows.

## ⚙️ Setup Requirements

To use this action, you need:

1. **Anthropic API Key**: Set your Claude API access key to `ANTHROPIC_API_KEY` secret
2. **Proper Permissions**: Workflows need permissions to modify PRs and issues
3. **Path Filter**: Currently configured to trigger only on `.github/**` file changes

## 🔍 Testing Purpose

This repository is used to test:

- Claude AI's PR description generation and formatting capabilities
- Claude AI's code review and inline commenting features
- Dual workflow orchestration (description + review)
- GitHub Actions workflow integration with Claude
- Path-based triggering for specific file changes
- Manual triggering via comment mentions
