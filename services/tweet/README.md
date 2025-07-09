# TweetService for Twittah

TweetService es un microservicio hexagonal responsable de la publicación de tweets y la gestión de timelines en la plataforma Twittah.

## 🚀 Visión
- Publicar y persistir tweets (Postgres/Ent).
- Fan-out a timelines cacheados en Redis de los followers.
- API REST para publicar y consultar timelines.
- Sincronización por eventos Kafka (follows, tweets).

## 🏗️ Arquitectura Hexagonal

```
services/tweet/
├── cmd/main.go
├── config/config.go
├── internal/
│   ├── domain/tweet_entity.go
│   ├── ports/tweet_repository.go
│   ├── usecase/tweet_usecase.go
│   └── adapters/
│       ├── http/tweet_handler.go
│       ├── kafka/tweet_consumer.go
│       ├── kafka/tweet_producer.go
│       ├── postgres/tweet_repository.go
│       └── redis/timeline_cache.go
```

- **Dominio**: entidad Tweet y validaciones.
- **Ports**: interfaces para persistencia y cache.
- **Usecase**: lógica de negocio orquestando dominio y puertos.
- **Adapters**: HTTP (Gin), Kafka, Postgres (Ent), Redis.

## 📦 Stack
- Go 1.21+
- Ent ORM
- PostgreSQL
- Redis (go-redis v8)
- Kafka (sarama)
- Gin
- Viper
- Zap

## 🌐 API REST
- `POST /tweets` — Publicar tweet
- `GET /timeline` — Obtener timeline de usuario

## ⚡ Eventos Kafka
- `tweets.published` (produce/consume)
- `follows.created` (consume)
- `follows.deleted` (consume)
- `timelines.updated` (produce)

## ⚙️ Configuración
Variables requeridas (ver `.env-example`):
- `POSTGRES_DSN`, `REDIS_ADDR`, `KAFKA_BROKERS`, `KAFKA_GROUP_ID`, `SERVICE_PORT`, `TIMELINE_TTL_HOURS`

## 🐳 Docker
Multi-stage build. Expone el puerto 8080.

## ▶️ Setup local
```sh
cp .env-example .env
# Completa tus variables si es necesario
make ent-generate # si usas Ent
make run
```

## 🧪 Testing
```sh
go test ./...
```

## 📝 Estándares
- Arquitectura hexagonal estricta
- Go idiomático y documentado (ver `.windsurf/rules/golang-standards.md`)
- Sin acceso directo a datos de usuario
- TTL de timelines en Redis = 24h

---

Desarrollado para Twittah siguiendo los más altos estándares de arquitectura y calidad Go.
