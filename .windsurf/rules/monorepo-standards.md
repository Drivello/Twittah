
# 🏛️ Monorepo Development Standards

## 🐹 Go Idioms and Best Practices

### Go Idioms and Best Practices
- Use `interface{}` sparingly, prefer explicit interfaces in ports and domain layers.
- Handle errors using error wrapping with context:
  ```go
  if err != nil {
      return fmt.Errorf("operation failed: %w", err)
  }
  ```
- Use `context.Context` in all public methods for cancellation and timeouts.
- Follow Hexagonal Architecture strictly: domain logic in usecase layer, infrastructure logic in adapters.
- Use dependency injection for repositories, cache, message brokers, and configuration.
- Implement graceful shutdown of HTTP servers, Kafka consumers/producers, and Redis/Postgres connections.
- Use structured logging with Zap or similar libraries.
- Keep all magic values and configs in `.env` files and load them with Viper or equivalent.
- Define constants for topics, keys, and error messages.
- For Goroutines:
  - Use `sync.WaitGroup` and `Context` for lifecycle management.
  - Avoid goroutine leaks in consumers or producers.

### Security
- Validate all incoming user inputs at the HTTP handler level.
- Sanitize and escape data before persisting or logging.
- Never log sensitive data (passwords, tokens, API keys).
- Apply the principle of least privilege in all external systems.
- Manage secrets with environment variables or vault systems.

### Performance
- Cache hot data with Redis or Memcached when appropriate.
- Use batched writes to cache or DB for high-throughput scenarios.
- Add proper indexes to DB queries to avoid full table scans.
- Use buffered channels for message producers to avoid blocking.
- Monitor and optimize resource usage in tight loops or consumers.

---

## 📝 Code Formatting Standards

### Line Spacing Rules
- Leave a blank line after closing instruction blocks (if, for, switch):
  ```go
  if condition {
      doSomething()
  }

  nextOperation()
  ```
- Leave a blank line before `return` statements.
- Combine assignment and error check when possible:
  ```go
  if err := repo.Save(ctx, entity); err != nil {
      return err
  }
  ```

---

## 📖 Documentation Standards

### GoDoc Documentation
- Document all exported types, functions, and interfaces.
- Focus on documenting:
  - Business logic (usecases)
  - Adapters that interact with external systems
  - Complex data flows or transformations
- Start comments with the function/type name.
- Include parameter and return value descriptions for complex methods.

Example:
```go
// ProcessEvent handles incoming Kafka events and updates the database.
func ProcessEvent(ctx context.Context, event Event) error {
    // implementation
}
```

---

## 🔄 Development Workflow
- Enforce pre-commit hooks for linting, formatting, and unit testing.
- PRs must be reviewed and approved before merging.
- Use semantic commit messages (Conventional Commits):
  ```
  feat(auth): add JWT token support
  fix(user): handle race condition in follow creation
  ```
- CI/CD pipelines should:
  - Build, lint, and test all affected services.
  - Push Docker images to a registry per service.

---

## 🐳 Docker and Deployment
- Multi-stage Docker builds for lightweight images.
- Root-level `docker-compose.yml` for local development:
  - Includes all services and dependencies (DB, cache, brokers).
- Each microservice exposes port `8080` internally.
- Healthchecks configured for all dependencies in Docker Compose.

---

## 📦 Caching and Message Brokers
### Redis Practices
- Use keys like `namespace:<resource_id>` for clarity.
- Set TTLs for ephemeral data (e.g., cache, sessions).
- Use pipelining or transactions for batch operations.
- Fail gracefully if cache is unavailable.

### Kafka (or other brokers)
- Use explicit topic names and partitioning strategies.
- Consumers should handle retries and dead-letter queues.
- Always acknowledge messages after successful processing.
- Prefer Sarama or similar Go clients for Kafka.

---

# 🏁 Summary
These standards promote consistency, scalability, and clean architecture in monorepos with multiple Go microservices. Apply them across all services to maintain a production-ready codebase.
