# User Service

Servicio encargado de la gestión de relaciones de seguimiento (follow/unfollow) entre usuarios.

## Funcionalidad
- Persiste una replica de usuarios en Postgres, los cuales consume a traves de eventos generados por Auth Service.
- Maneja operaciones de follow y unfollow vía consumidores de Kafka, generados en Gateway Service. 
- Expone endpoints para consultar followers y followees.
- Pensado para escalar y desacoplar las lecturas de las escrituras.

## Stack
- Go 1.23
- Gin
- Kafka (Sarama)
- Ent (Postgres)
- Redis
- Uber Zap (logging)
- Docker

## Variables de entorno principales
Ver `.env-example` para la lista completa y ejemplos.
- `USER_SERVICE_PORT`: Puerto HTTP
- `LOG_LEVEL`: Nivel de logs
- `USER_POSTGRES_DSN`: DSN de Postgres
- `REDIS_ADDR`: Redis
- `KAFKA_BROKERS`: Brokers de Kafka

## Docker
Se recomienda levantar junto al resto de microservicios usando `docker-compose`.

## Notas
- El servicio es event-driven y desacoplado.
- Futuro: perfil de usuario, bio, avatar, etc.
