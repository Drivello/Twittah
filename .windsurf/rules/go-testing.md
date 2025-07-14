---
trigger: manual
---


# Unit Testing Rules

## Testing Best Practices

### Testing Principles
- Write tests that are easy to read, understand, and maintain
- Test behavior, not implementation details
- Follow the AAA pattern (Arrange, Act, Assert)
- Use descriptive test names that explain the scenario and expected outcome
- Test both success and error cases
- Test edge cases and error conditions

### Go Testing Patterns

#### Table-Driven Tests
Follow the table-driven test pattern for comprehensive test coverage:

```go
// TEMPLATE: Adapt this example to your specific use case
// 1. Replace "Function" with the name of the function you want to test
// 2. Define the "args" struct with your function's input parameters
// 3. Define the "wants" struct with the expected output values
// 4. Add test cases to the "tests" slice covering success, error, and edge case scenarios
// 5. Adjust the function call and assertions according to your business logic
func TestFunction(t *testing.T) {
    tests := []struct {
        name  string
        args  args
        wants wants
    }{
        {
            name: "should_return_success_when_valid_input",
            args: args{
                input: "valid",
            },
            wants: wants{
                result: "expected",
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := Function(test.args.input)
            assert.Equal(t, test.wants.result, result)
        })
    }
}
```

#### Test Structure and Organization
- Use descriptive test names that explain the scenario and expected outcome
- Structure test cases with clear input (args) and expected output (wants)
- Use testify/assert for assertions
- Mock external dependencies using interfaces
- Test both success and error cases
- Use subtests for better organization and selective test execution
- Test edge cases and error conditions
- Use test fixtures for complex test data
- Follow AAA pattern (Arrange, Act, Assert) in test cases

### Testing Best Practices
- Test behavior, not implementation
- Keep tests simple and focused
- Use meaningful test names
- Test edge cases and error conditions
- Mock external dependencies
- Use table-driven tests for multiple scenarios
- Organize tests logically
- Maintain test data separately
- Use setup and teardown appropriately
- Keep tests independent of each other