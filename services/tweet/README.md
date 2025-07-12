# Microservicio Tweet

Este servicio gestiona los tweets y el timeline de usuarios.

## Funcionalidades principales
- Replica mínima de usuarios (consume `users.created` de Kafka)
- Manejo de eventos de tweets (`tweet.created`, `tweet.deleted`)
- Endpoint REST GET /timeline
- Caching con Redis
- Observabilidad (Prometheus, Zap)
- Healthchecks

## Stack
- Go 1.23
- Ent (Postgres)
- Redis
- Kafka (Sarama)
- Zap
- Prometheus
- Docker

## Uso

```bash
# Build
make build

# Run
make run

# Test
make test
```

## Docker
Levantar con Docker Compose junto a los otros servicios.
