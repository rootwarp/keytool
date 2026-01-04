---
name: go-review
description: Review Go code for best practices, security issues, and performance. Use for code reviews or quality checks.
allowed-tools: Read, Grep, Glob
---

# Go Code Review Skill

Review Go code for quality, security, and adherence to project standards.

## Review Checklist

### Code Quality

- [ ] Functions are small and focused (single responsibility)
- [ ] Variable names are descriptive and meaningful
- [ ] No magic numbers or strings (use constants)
- [ ] No dead code or commented-out code
- [ ] Code is formatted with `gofmt`

### Error Handling

- [ ] All errors are handled (not ignored with `_`)
- [ ] Errors are wrapped with context: `fmt.Errorf("context: %w", err)`
- [ ] No panics in library code (only in main/init if necessary)
- [ ] Error messages are lowercase and don't end with punctuation

### Documentation

- [ ] Exported functions have Godoc comments
- [ ] Comments explain "why", not "what"
- [ ] Package has a doc comment

### Interfaces

- [ ] Interfaces are defined in consumer packages, not providers
- [ ] Interfaces are small and focused
- [ ] Dependency injection is used for testability

### Testing

- [ ] New code has corresponding tests
- [ ] Tests cover both success and error paths
- [ ] Table-driven tests used for multiple cases
- [ ] Test names clearly describe what's being tested

### Performance

- [ ] No unnecessary allocations in hot paths
- [ ] Slices pre-allocated when size is known
- [ ] Context used for cancellation and timeouts
- [ ] No goroutine leaks (all goroutines properly terminated)

## Security Review (Ethereum-Specific)

### Key Management

- [ ] Key material is never logged
- [ ] Keys are properly encrypted at rest
- [ ] Memory is zeroed after key use when possible
- [ ] No hardcoded keys or mnemonics

### Input Validation

- [ ] Ethereum addresses are validated before use
- [ ] Public keys are validated (correct length, on curve)
- [ ] Amounts and values are bounds-checked
- [ ] Network IDs are validated

### Secret Management

- [ ] Secrets use Google Secret Manager
- [ ] No credentials in code or config files
- [ ] Environment variables used appropriately

### BLS/Cryptography

- [ ] Signatures are verified before trusting data
- [ ] Domain separation is correct
- [ ] Fork version is properly handled

## Common Issues to Flag

### Anti-Patterns

```go
// Bad: Ignoring error
result, _ := doSomething()

// Bad: Panic in library code
if err != nil {
    panic(err)
}

// Bad: Error not wrapped
return err

// Bad: Empty interface overuse
func process(data interface{}) {}
```

### Better Alternatives

```go
// Good: Handle error
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}

// Good: Return error in library code
if err != nil {
    return fmt.Errorf("unexpected condition: %w", err)
}

// Good: Error wrapped with context
return fmt.Errorf("processing request: %w", err)

// Good: Specific types
func process(data *RequestData) error {}
```

## Review Output Format

When reviewing code, provide:

1. **Summary**: Overall assessment (Approved/Changes Requested)
2. **Critical Issues**: Security vulnerabilities, bugs, data loss risks
3. **Suggestions**: Improvements, better patterns, optimizations
4. **Nits**: Minor style or formatting issues

Use file:line references for all comments.
