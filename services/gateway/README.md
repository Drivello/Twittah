# Gateway Service

Este servicio actúa como el API Gateway y orquestador principal de la plataforma tipo Twitter.

## Funcionalidad
- Expone la API pública REST para clientes.
- Orquesta llamadas a UserService y TweetService para operaciones de lectura.
- Aplica middleware de autenticación real (consulta a AuthService en cada request).
- Propaga eventos a través de Kafka para operaciones de escritura de autenticacion, usuario y tweets.
- Healthcheck.

## Stack
- Go 1.23
- Gin
- Kafka (Sarama)
- Uber Zap (logging)
- Docker

## Variables de entorno principales
Ver `.env-example` para la lista completa y ejemplos.
- `GATEWAY_PORT`: Puerto HTTP
- `LOG_LEVEL`: Nivel de logs
- `KAFKA_BROKERS`: Brokers de Kafka
- `KAFKA_USER_TOPIC`, `KAFKA_AUTH_TOPIC`, `KAFKA_TWEET_TOPIC`: Topics usados
- `AUTH_MICROSERVICE_URL`: URL del AuthService
- `USER_MICROSERVICE_URL`: URL del UserService
- `TWEET_MICROSERVICE_URL`: URL del TweetService

## Docker
Se recomienda levantar junto al resto de microservicios usando `docker-compose`.

## Arquitectura
- Hexagonal (Ports & Adapters)
- Middleware global de autenticación: valida cada request con AuthService usando el header `X-User-Id`.
- Pensado para desacoplar lógica de negocio y facilitar testing.

