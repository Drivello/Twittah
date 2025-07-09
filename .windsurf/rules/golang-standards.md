---
trigger: manual
---


# Go Development Standards

## Go Idioms and Best Practices

### Go Idioms and Best Practices
- Use `interface{}` sparingly, prefer specific interfaces
- Implement error handling using custom error types
- Use `context.Context` for cancellation and timeouts
- Follow standard Go project layout
- Use dependency injection for better testability
- Implement graceful shutdown
- Use structured logging
- Follow Go's error handling patterns:
  ```go
  if err != nil {
      return fmt.Errorf("failed to process user: %w", err)
  }
  ```

### Security
- Validate all user inputs
- Implement authentication and authorization where needed
- Protect sensitive data
- Follow the principle of least privilege

### Performance
- Optimize database queries
- Implement caching when appropriate
- Monitor performance
- Handle resources efficiently
- Use goroutines safely, avoiding infinite loops and race conditions
- Prefer buffered channels over unbuffered channels for better performance

## Code Formatting Standards

### Line Spacing Rules
- **Leave a blank line after closing instruction blocks** (conditionals, loops, etc.)
  ```go
  // Good
  if condition {
      doSomething()
  }

  nextStatement()

  // Bad
  if condition {
      doSomething()
  }
  nextStatement()
  ```

- **Leave a blank line before `return` statements**
  ```go
  // Good
  func processData() error {
      data := getData()      
      if data == nil {
          return errors.New("no data")
      }

      return nil
  }

  // Bad
  func processData() error {
      data := getData()
      
      if data == nil {
          return errors.New("no data")
      }
      return nil
  }
  ```

- **Leave a blank line before `if` statements, except when the previous line defines a variable used in the `if` condition**
  ```go
  // Good - variable used in if condition
  err, data := processData()
  if err != nil {
      return err
  }

  // Good - blank line before if
  processData()

  if someCondition {
      doSomething()
  }

  // Bad - unnecessary blank line when variable is used in if
  err := processData()

  if err != nil {
      return err
  }
  ```

- **Prefer inline error checking when possible**
  ```go
  // Good - combine assignment and error check
  if err := ctrl.Get(ctx, categorizationID, siteID, &category); err != nil {
      return categorydto.CategoryResponse{}, err
  }

  // Less preferred - separate assignment and check
  err = ctrl.Get(ctx, categorizationID, siteID, &category)
  if err != nil {
      return categorydto.CategoryResponse{}, err
  }
  ```

## Documentation Standards

### GoDoc Documentation
- Document all exported functions, methods, types, and constants using godoc format
- Focus on documenting functions and methods of significant importance or high complexity
- Write clear, concise documentation that explains the purpose, parameters, and return values
- Use proper godoc formatting with the function/method name as the first word
- Document complex business logic and non-obvious behavior
- Include examples for complex or frequently used functions

#### Functions of Significant Importance or High Complexity
Prioritize documentation for functions that meet any of these criteria:
- **Business Logic Functions**: Core business operations, calculations, or domain-specific algorithms
- **Complex Processing**: Functions with multiple steps, conditional logic, or data transformations
- **Error-Prone Operations**: Functions that handle external dependencies, network calls, or file operations
- **Public APIs**: Functions that serve as entry points or interfaces for other modules
- **Performance-Critical**: Functions that handle large datasets or have performance implications
- **Non-Obvious Behavior**: Functions with side effects, special return conditions, or edge case handling
- **Configuration/Setup**: Functions that initialize systems, configure services, or manage application state

Examples:
```go
// ProcessUserData validates and processes user information for storage.
// It returns the processed user data and any validation errors encountered.
// The function handles email normalization and password hashing automatically.
func ProcessUserData(userData *UserInput) (*User, error) {
    // implementation
}

// CalculateMonthlySpending computes the total spending for a user within a specific month.
// It aggregates all transactions for the given userID and month, applying any necessary
// currency conversions and filtering out refunds.
//
// Parameters:
//   - userID: unique identifier for the user
//   - month: target month in YYYY-MM format
//   - includePending: whether to include pending transactions
//
// Returns the total spending amount and any computation errors.
func CalculateMonthlySpending(userID string, month string, includePending bool) (float64, error) {
    // implementation
}
```

## Development Workflow
- Perform code reviews before merging
- Keep code updated with main branch
- Write descriptive commits
- Document significant changes
- Maintain an updated changelog