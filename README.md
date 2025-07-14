# 🐦 Twittah – Microservices Monorepo

**Twittah** is a microblogging system inspired by Twitter. This monorepo contains all microservices needed for user creation, tweet publishing, user following relationships, and optimized timelines using Redis cache.

> 🚀 Designed to scale to millions of users with a focus on clean architecture, resilience, and observability.

---

## 📂 Monorepo Contents

| Service        | Description                                          |
|----------------|------------------------------------------------------|
| `auth`         | User authentication and validation                   |
| `user`         | Manages users and follow/unfollow relationships      |
| `tweet`        | Manages tweets and timelines with caching            |
| `gateway`      | Orchestrates calls between microservices             |

---

## 🏛️ Architecture

- **Hexagonal (Ports & Adapters):** Each service is decoupled and highly testable.
- **Kafka:** Event bus for asynchronous communication.
- **PostgreSQL:** Independent databases per service.
- **Redis:** Timeline caching in `tweet` for faster response times.
- **Docker Compose:** Full local development environment (👉 **RECOMMENDED: use Docker Compose**).
- **Go 1.23:** Main language for all services.
- **Ent ORM:** Strongly typed data models with safe migrations. 💥

```
                           +-----------+
                           |  Gateway  |
                           +-----------+
                                |
         +---------------------+---------------------+
         |                     |                     |
    +--------+            +--------+            +--------+
    |  Auth  |            |  User  |            |  Tweet  |
    +--------+            +--------+            +--------+
         |                     |                     |
  PostgreSQL            PostgreSQL            PostgreSQL + Redis
```

---

## 🚀 Quick Start

### 1️⃣ Prerequisites
- Docker & Docker Compose (**Essential**, all setup is designed for Compose)
- Go 1.23

### 2️⃣ Clone the repository
```bash
git clone https://github.com/Drivello/Twittah.git
cd twittah
```

### 3️⃣ Start the full environment

Make sure to have `make` installed. You can get it with:
```bash
choco install make
```

Run:
```bash
make all-up
```

👉 **Note:** The `docker-compose.yml` automatically handles dependencies:
- Kafka
- PostgreSQL
- Redis
- Grafana
- Prometheus


### 4️⃣ Access Postman Collections

Use this Postman Collection: [Postman Collection](Twittah_Gateway.postman_collection.json)

### 5️⃣ Access Grafana (metrics)

Go to: http://localhost:3000

- user: admin
- password: admin

Create a new Prometheus Data Source at: 
http://localhost:3000/connections/datasources

Configure the Data Source: 
- Name: Prometheus
- URL: http://prometheus:9090

Save & Test → Build Dashboard → Import Dashboard

Drag and drop the dashboard JSON file

[Dashboard.json](monitoring/grafana/twittah-dashboard.json)

Enjoy!

---

## 🧪 Testing

### Current status
✅ Unit tests implemented in `gateway`.  
⚠️ Other services have no test coverage yet (marked as a priority for upcoming iterations).

### Run existing tests
```bash
cd services/gateway
go test -cover -count=1 ./...
```

---

## 📝 Logging

The project uses **structured logging with Uber Zap**.  

---

## 📦 Main Directory Structure

```
services/
│
├── auth/       # Authentication microservice 
├── user/       # User management microservice
├── tweet/      # Tweets and timeline microservice
├── gateway/    # API orchestrator
│
docker-compose.yml
```

---

## 🛠️ Key Technologies

| Component          | Technology               |
|--------------------|---------------------------|
| Backend            | Go 1.23                   |
| Communication      | Kafka (event-driven)      |
| Persistence        | PostgreSQL + Ent ORM      |
| Cache              | Redis                     |
| Containerization   | Docker & Docker Compose   |
| Configuration      | Viper                     |

---

## 📦 Simplified Data Flow

1. User sends request to `Gateway`.
2. `Gateway` validates and orchestrates calls to other services.
3. Important events are sent via Kafka.
4. Services react and update their state (DB or cache).

---

## 🏁 Roadmap (Post-MVP)

### ✅ Quality and Testing
- Add test coverage to all microservices (unit >70%, e2e for critical flows).
- Improve existing tests in Gateway for edge cases.
- Better HTTP error responses.
- Add load and stress tests for scalability validation.
- Implement circuit breakers for external services.

### 🔐 Security
- Add Login/Signup system. 
- Add authentication with Bearer Tokens.

### 🌐 API Gateway
  - Introduce Kong in front of Gateway for rapid scaling.
  - Per-user/IP rate limiting.
  - Centralized access logging.
  - Microservices in private networks, Gateway as the only exposed entry point.
  - CORS validation.

### 📊 Observability
- Add alerts (high latency, 5xx errors, Kafka/Redis failures).
- Enriched structured logs with tracing IDs.

### 🚀 Scalability and Optimization
- Full fan-out on write in TweetService to distribute tweets to cached timelines.
- Cache paginated timelines in Redis (`timeline:<user_id>:page:<n>`).
- Use prefetch and pipelining in Redis to reduce cache-miss latency.
- Kafka payload compression (Snappy or LZ4).
- Consider PostgreSQL sharding in TweetService if required.

### ☁️ Infrastructure and Deployment
- Deploy to AWS ECS or EKS using Helm charts.
- Infrastructure as Code (Terraform).
- Full CI/CD pipeline with GitHub Actions or GitLab CI (build, lint, test, deploy to staging/prod).

### 👨‍💻 Developer Experience
- Add pre-commit hooks (formatting, linting, tests).
- Detailed documentation per microservice.

---

## 📬 Contact

**Nicolás Sánchez** – [LinkedIn](https://www.linkedin.com/in/mario-nahuel-nicolas-sanchez)  
📧 nico_dd@outlook.com.ar