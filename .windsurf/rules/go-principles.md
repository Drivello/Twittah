---
trigger: manual
---


# Go-Specific Examples for Development Principles

> For general development principles (KISS, YAGNI, DRY, etc.), see [principles/rule.mdc](mdc:../principles/rule.mdc)  
> For software architecture principles (SOLID, Separation of Concerns, etc.), see [architecture/rule.mdc](mdc:../architecture/rule.mdc)

This file contains Go-specific examples and implementations of the general development principles.

## Naming Conventions in Go

### Function Names
```go
// Good - Intention revealing with Go conventions
func isValidDateFormat() bool
func hasPermission() bool
func canAccess() bool
func validateUser() error

// Bad - Unclear intention or non-Go conventions  
func checkDateFormat() bool  // what does check mean? validation? verification?
func userValidation() error  // noun instead of verb
func queryParamDateFormat() bool  // implementation detail in name
```

### Go-Specific Naming Rules
- Package names: short, lowercase, single words
- Exported functions: PascalCase (`ValidateUser`)
- Unexported functions: camelCase (`validateUser`)
- Constants: mixed caps (`MaxRetries`) or all caps for exported (`MAX_RETRIES`)
- Interfaces: often end with `-er` (`Reader`, `Writer`, `Validator`)

## Return Early Pattern in Go

```go
// Good - Early returns reduce nesting
func processUser(userID string) (*User, error) {
    if userID == "" {
        return nil, errors.New("userID cannot be empty")
    }

    user, err := fetchUser(userID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch user: %w", err)
    }

    if !user.IsActive {
        return nil, errors.New("user is not active")
    }

    return user, nil
}

// Bad - Deep nesting makes code harder to read
func processUser(userID string) (*User, error) {
    if userID != "" {
        user, err := fetchUser(userID)
        if err == nil {
            if user.IsActive {
                return user, nil
            } else {
                return nil, errors.New("user is not active")
            }
        } else {
            return nil, fmt.Errorf("failed to fetch user: %w", err)
        }
    } else {
        return nil, errors.New("userID cannot be empty")
    }
}
```

## Composition over Inheritance in Go

Go doesn't have traditional inheritance but uses embedding and interfaces:

```go
// Good - Composition with embedding
type User struct {
    Person    // Embedded struct
    Role   string
    Active bool
}

type Person struct {
    Name  string
    Email string
}

// Interface composition
type ReadWriter interface {
    Reader
    Writer
}

// Bad - Trying to mimic inheritance incorrectly
type User struct {
    PersonName  string
    PersonEmail string
    Role        string
    Active      bool
}
```

## DRY in Go

```go
// Good - Extract common validation logic
func validateEmail(email string) error {
    if email == "" {
        return errors.New("email cannot be empty")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    return nil
}

func CreateUser(name, email string) error {
    if err := validateEmail(email); err != nil {
        return fmt.Errorf("user creation failed: %w", err)
    }
    // Create user logic
    return nil
}

func UpdateUser(id int, email string) error {
    if err := validateEmail(email); err != nil {
        return fmt.Errorf("user update failed: %w", err)
    }
    // Update user logic
    return nil
}

// Bad - Repeated validation logic
func CreateUser(name, email string) error {
    if email == "" {
        return errors.New("email cannot be empty")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    // Create user logic
    return nil
}

func UpdateUser(id int, email string) error {
    if email == "" {
        return errors.New("email cannot be empty")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email format")
    }
    // Update user logic
    return nil
}
```

## Fail Fast in Go

```go
// Good - Validate early and fail fast
func ProcessPayment(amount float64, currency string, userID string) error {
    // Validate inputs immediately
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    if currency == "" {
        return errors.New("currency is required")
    }
    if userID == "" {
        return errors.New("userID is required")
    }
    
    // Main business logic after validation
    return processPaymentInternal(amount, currency, userID)
}

// Bad - Validation mixed with business logic
func ProcessPayment(amount float64, currency string, userID string) error {
    user, err := fetchUser(userID)
    if err != nil {
        return err
    }
    
    payment := calculatePayment(amount, currency)
    if amount <= 0 {  // Too late - already did expensive operations
        return errors.New("amount must be positive")
    }
    
    return savePayment(payment)
}