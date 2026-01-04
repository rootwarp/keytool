---
name: go-testing
description: Write Go unit tests following project patterns. Use when creating tests, improving coverage, or fixing test failures.
allowed-tools: Read, Grep, Glob, Edit, Bash(go test:*)
---

# Go Testing Skill

Write Go tests following this project's conventions and best practices.

## Test File Location

- Tests go in `*_test.go` files alongside the code they test
- Use the same package name for white-box testing
- Use `_test` package suffix for black-box testing

## Test Patterns

### Table-Driven Tests (Preferred)

```go
func TestFunctionName_Scenario(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "valid input returns expected output",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "invalid input returns error",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tt.want, got)
        })
    }
}
```

### Assertions

Use `testify/require` for assertions:

```go
import "github.com/stretchr/testify/require"

require.NoError(t, err)
require.Error(t, err)
require.Equal(t, expected, actual)
require.NotNil(t, obj)
require.True(t, condition)
require.Contains(t, str, substring)
```

## Mocking

This project uses mockery for mock generation. Configuration is in `.mockery.yaml`.

### Using Mocks

```go
import "github.com/dsrv/eth-staking-key-operator/mocks"

func TestWithMock(t *testing.T) {
    mockWallet := mocks.NewMockWallet(t)
    mockWallet.EXPECT().GetKey(ctx, keyID).Return(key, nil)

    svc := NewService(mockWallet)
    result, err := svc.DoSomething(ctx)

    require.NoError(t, err)
}
```

### Generating New Mocks

If you add a new interface that needs mocking:
1. Add it to `.mockery.yaml`
2. Run `mockery` to generate

## Test Coverage

- Target >80% coverage for new code
- Test both success and error paths
- Test edge cases (nil inputs, empty slices, boundary values)

## Running Tests

```bash
go test -v ./...                    # All tests verbose
go test -v -run TestName ./...      # Single test
go test -cover ./...                # With coverage
go test -race ./...                 # Race detection
make test                           # Project's test target
```

## Common Test Scenarios

### Testing Error Cases

```go
{
    name:    "nil context returns error",
    ctx:     nil,
    wantErr: true,
},
{
    name:    "empty input returns error",
    input:   "",
    wantErr: true,
},
```

### Testing with Context

```go
func TestWithContext(t *testing.T) {
    ctx := context.Background()
    // or with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
}
```
