# Tweet Service

Servicio encargado de la gestión de tweets y timelines de usuarios para la plataforma tipo Twitter.

## Funcionalidad
- Publicar y eliminar tweets vía consumidores Kafka (`tweets.published`, `tweets.deleted`) generados en Gateway Service.
- Consume eventos `users.created` para replicar usuarios en Postgres.
- Endpoint REST `/timeline` para obtener el timeline del usuario.
- Cache de timelines en Redis (TTL configurable).
- Observabilidad (Zap).
- Healthcheck.

## Stack
- Go 1.23
- Gin
- Ent (Postgres)
- Redis
- Kafka (Sarama)
- Uber Zap (logging)
- Docker

## Variables de entorno principales
Ver `.env-example` para la lista completa y ejemplos.
- `TWEET_PORT`: Puerto HTTP
- `LOG_LEVEL`: Nivel de logs
- `TWEET_POSTGRES_DSN`: DSN de Postgres
- `TWEET_REDIS_ADDR`: Redis
- `KAFKA_BROKERS`: Brokers de Kafka
- `KAFKA_TWEET_TOPIC`, `KAFKA_TWEET_DLQ_TOPIC`, ...: Topics usados

## Docker
Se recomienda levantar junto al resto de microservicios usando `docker-compose`.

## Notas
- El timeline se cachea por usuario con TTL (ver variable en `.env-example`).
- Consumir eventos `follows.created`, `follows.deleted` para mantener consistencia de timelines y usuarios.
- Futuro: fan-out on write, paginación avanzada, soporte para multimedia.
