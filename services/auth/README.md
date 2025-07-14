# Auth Service

Servicio de autenticación y validación de usuarios para la plataforma tipo Twitter.

## Funcionalidad
- Valida el usuario recibido en el header `X-User-Id` (por ahora solo que sea un entero positivo).
- Provee endpoint `/validate` para validación de usuario (usado por Gateway).
- Pensado para evolucionar a login/signup y JWT en el futuro.
- Healthcheck.

## Stack
- Go 1.23
- Gin
- Kafka (Sarama)
- Ent (Postgres)
- Uber Zap (logging)
- Docker

## Variables de entorno principales
Ver `.env-example` para la lista completa y ejemplos.
- `AUTH_PORT`: Puerto HTTP
- `LOG_LEVEL`: Nivel de logs
- `AUTH_POSTGRES_DSN`: DSN de Postgres
- `KAFKA_BROKERS`: Brokers de Kafka
- `KAFKA_USER_TOPIC`, `KAFKA_AUTH_TOPIC`: Topics usados

## Docker
Se recomienda levantar junto al resto de microservicios usando `docker-compose`.

## Notas
- El endpoint `/validate` es un mock, pero el middleware en Gateway es real.
- Futuro: login/signup, JWT, OAuth2.
