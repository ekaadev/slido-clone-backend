# GitHub Actions Workflows

## Code Review Workflows

### Primary: Claude Code Review
**File:** `claude-code-review.yml`

Automatically reviews PRs using Claude AI. This is the primary code review workflow.

- **Trigger:** Runs on PR opened, synchronize, ready_for_review, reopened
- **Status:** Non-blocking (uses `continue-on-error: true`)
- **Behavior:** If Claude hits rate limits, the workflow will fail but won't block the PR

### Backup: AI Code Review
**File:** `copilot-code-review.yml`

Backup code review using OpenAI GPT-4 when Claude is unavailable.

- **Trigger:** 
  - Automatically on PR events (same as Claude)
  - Manually via workflow_dispatch
- **Requires:** `OPENAI_API_KEY` secret
- **Behavior:** Falls back to manual review notice if all AI services fail

#### Manual Trigger
When Claude hits rate limits, you can manually trigger the backup review:

```bash
# Via GitHub CLI
gh workflow run copilot-code-review.yml -f pr_number=13

# Via GitHub UI
# Go to Actions → AI Code Review (Backup) → Run workflow
```

## Setup

### Required Secrets

1. **CLAUDE_CODE_OAUTH_TOKEN** (already configured)
   - For Claude code review workflow

2. **OPENAI_API_KEY** (optional, for backup)
   - Get from: https://platform.openai.com/api-keys
   - Add in: Repository Settings → Secrets → Actions
   - Format: `sk-...` (your OpenAI API key)

### Configuration

To use the backup workflow, add the `OPENAI_API_KEY` secret:

```bash
# Using GitHub CLI
gh secret set OPENAI_API_KEY --body "sk-your-api-key-here"
```

## Workflow Logic

```
PR opened/updated
    ↓
Claude Review (non-blocking)
    ↓
    ├─ Success → Post review
    ├─ Failure (rate limit) → PR not blocked
    └─ If needed → Manually trigger backup workflow
```

## Current Status

- ✅ Claude workflow is non-blocking
- ✅ Backup workflow ready (requires OPENAI_API_KEY)
- ✅ Manual trigger available for backup
- ✅ Fallback to manual review if all AI fails

## Troubleshooting

### Claude hits rate limit
1. The PR is not blocked (workflow marked as non-blocking)
2. Manually trigger backup: `gh workflow run copilot-code-review.yml -f pr_number=<PR#>`
3. Or perform manual code review

### Backup also fails
- Check if OPENAI_API_KEY is configured
- Verify the secret is valid
- Fall back to manual code review
