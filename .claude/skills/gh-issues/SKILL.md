---
name: gh-issues
description: Manage GitHub issues with full CRUD operations. Use for creating, listing, updating, commenting on, and closing issues.
allowed-tools: Bash(gh issue:*), Bash(gh api:*), Read
---

# GitHub Issues Skill

Manage GitHub issues using the `gh` CLI. Supports full CRUD operations including creating, listing, searching, updating, commenting, and closing issues.

## Prerequisites

Ensure `gh` is authenticated:
```bash
gh auth status
```

For project operations, add project scope:
```bash
gh auth refresh -s project
```

## Create Issues

### Basic Creation
```bash
gh issue create --title "Bug: Application crashes" --body "Steps to reproduce..."
```

### With Labels and Assignees
```bash
gh issue create \
  --title "Feature: Add dark mode" \
  --body "Users have requested dark mode support" \
  --label "enhancement,ui" \
  --assignee "@me"
```

### From File (for longer descriptions)
```bash
gh issue create --title "Detailed bug report" --body-file ./bug-description.md
```

### In Another Repository
```bash
gh issue create -R owner/repo --title "Cross-repo issue" --body "Details"
```

## List and Search Issues

### Basic Listing
```bash
gh issue list                              # Open issues (default)
gh issue list --state closed               # Closed issues
gh issue list --state all                  # All issues
```

### Filtering
```bash
gh issue list --label "bug"                # By label
gh issue list --label "bug,critical"       # Multiple labels
gh issue list --assignee "@me"             # Assigned to me
gh issue list --author "username"          # By author
gh issue list --milestone "v1.0"           # By milestone
```

### Advanced Search
```bash
gh issue list --search "error no:assignee sort:created-asc"
gh issue list --search "is:open label:bug created:>2024-01-01"
```

### JSON Output for Processing
```bash
gh issue list --json number,title,labels,state
gh issue list --json number,title --jq '.[] | "\(.number): \(.title)"'
```

## View Issue Details

```bash
gh issue view 123                          # View issue
gh issue view 123 --comments               # Include comments
gh issue view 123 --web                    # Open in browser
gh issue view 123 --json body,comments     # JSON output
```

## Comment on Issues

```bash
gh issue comment 123 --body "Investigation update: found the root cause"
gh issue comment 123 --body-file ./detailed-analysis.md
gh issue comment 123 --edit-last --body "Updated comment"
gh issue comment 123 --delete-last --yes
```

## Update Issues

### Modify Content
```bash
gh issue edit 123 --title "Updated title"
gh issue edit 123 --body "New description"
gh issue edit 123 --body-file ./new-description.md
```

### Manage Labels
```bash
gh issue edit 123 --add-label "in-progress,priority:high"
gh issue edit 123 --remove-label "needs-triage"
```

### Manage Assignees
```bash
gh issue edit 123 --add-assignee "@me"
gh issue edit 123 --add-assignee "user1,user2"
gh issue edit 123 --remove-assignee "user1"
```

### Manage Milestones and Projects
```bash
gh issue edit 123 --milestone "v2.0"
gh issue edit 123 --remove-milestone
gh issue edit 123 --add-project "Roadmap"
```

### Batch Updates
```bash
gh issue edit 100 101 102 --add-label "sprint-5"
```

## Close and Reopen Issues

### Close
```bash
gh issue close 123
gh issue close 123 --reason "completed"
gh issue close 123 --reason "not planned"
gh issue close 123 --comment "Fixed in PR #456"
```

### Reopen
```bash
gh issue reopen 123
gh issue reopen 123 --comment "Reopening due to regression"
```

## Common Workflows

### Bug Triage
```bash
# Find unassigned bugs
gh issue list --label "bug" --search "no:assignee"

# Assign and prioritize
gh issue edit 123 --add-assignee "@me" --add-label "priority:high"

# Add triage note
gh issue comment 123 --body "Triaged: Regression from v1.2, investigating"
```

### Sprint Planning
```bash
# List milestone issues
gh issue list --milestone "Sprint 5" --state all

# Batch assign to sprint
gh issue edit 200 201 202 --add-label "current-sprint"
```

## Error Handling

| Error | Resolution |
|-------|------------|
| `Not logged in` | Run `gh auth login` |
| `Could not resolve to a Repository` | Check repo name: `gh repo view` |
| `Could not resolve to a Project` | Add scope: `gh auth refresh -s project` |
| `issue not found` | Verify: `gh issue list` |
| `HTTP 403` | Check repository permissions |
| `HTTP 422: Validation Failed` | Verify labels exist: `gh label list` |

### Diagnostic Commands
```bash
gh auth status                             # Check authentication
gh repo view                               # Verify repository
gh label list                              # List available labels
```

## Best Practices

1. **Clear titles**: Use prefixes (`Bug:`, `Feature:`, `Docs:`) for categorization
2. **Consistent labels**: Follow naming conventions (e.g., `type:bug`, `priority:high`)
3. **Focused assignments**: Limit to 1-2 assignees per issue
4. **Comment for updates**: Preserve history by commenting rather than editing body
5. **Closing comments**: Always explain resolution when closing
6. **Cross-references**: Link related issues with `#123` syntax
