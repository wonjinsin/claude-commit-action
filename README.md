# Claude Commit Action Test Repository

This repository is created for testing the Claude PR Assistant GitHub Action.

## 📁 .github Directory Structure

```
.github/
├── PULL_REQUEST_TEMPLATE.md
└── workflows/
    ├── claude.yml
    └── root_workflows.yml
```

## 🔧 GitHub Actions Workflows

### 1. claude.yml - Claude PR Assistant Workflow

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

### 2. root_workflows.yml - Usage Example

A workflow that actually calls the Claude PR Assistant.

**Trigger Conditions:**

- When a PR is opened (`pull_request.opened`)
- When `@claude_pr` is mentioned in issue comments

**Permissions:**

- `contents: read` - Read repository contents
- `pull-requests: write` - Modify PRs
- `issues: write` - Write issues
- `id-token: write` - Write ID tokens

## 🚀 Usage

1. **On PR Creation**: When creating a new PR, Claude automatically analyzes changes and fills the template.

2. **Manual Trigger**: Mention `@claude_pr` in PR comments to re-run Claude.

## ⚙️ Setup Requirements

To use this action, you need:

1. **Anthropic API Key**: Set your Claude API access key to `ANTHROPIC_API_KEY` secret
2. **Proper Permissions**: Workflow needs permissions to modify PRs and issues

## 🔍 Testing Purpose

This repository is used to test:

- Claude AI's PR analysis and description generation capabilities
- Proper operation of GitHub Actions workflows
- Automatic Pull Request template filling functionality
- Behavior under various trigger conditions
