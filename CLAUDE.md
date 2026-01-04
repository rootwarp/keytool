# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

keytool is a Go project. Module: `github.com/rootwarp/keytool`

## Build Commands

```bash
make build           # Build binary to bin/keytool
make test            # Run all tests
make coverage        # Generate coverage report (coverage.html)
make clean           # Remove build artifacts

# Run a specific test
go test -run TestName ./path/to/package
```

## Project Structure

Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

- `cmd/keytool/` - Main application entry point
- `internal/` - Private application code (not importable by other projects)
- `pkg/` - Library code that can be used by external applications
- `docs/` - Documentation and design documents
- `scripts/` - Build, install, and analysis scripts

## Code Conventions

- Maximum line length: 100 columns
- Follow SOLID principles:
  - **S**ingle Responsibility: Each module/function has one reason to change
  - **O**pen/Closed: Open for extension, closed for modification
  - **L**iskov Substitution: Subtypes must be substitutable for their base types
  - **I**nterface Segregation: Many specific interfaces over one general-purpose interface
  - **D**ependency Inversion: Depend on abstractions, not concretions

## Development Workflow (TDD)

Follow Kent Beck's Red-Green-Refactor cycle:

1. **Red**: Write a single failing test for the new feature
2. **Green**: Write minimal code to make the test pass
3. **Refactor**: Improve the code while keeping tests green

Repeat for each new feature or behavior.
