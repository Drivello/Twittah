
# 🐦 Twittah – Microservices Monorepo

**Twittah** es un sistema de microblogging inspirado en Twitter. Este monorepo contiene todos los microservicios necesarios para permitir la creación de usuarios, publicación de tweets, seguimiento entre usuarios y la generación de timelines optimizados con Redis cache.

> 🚀 Diseñado para escalar a millones de usuarios con un enfoque en arquitectura limpia, resiliencia y observabilidad.

---

## 📂 Contenido del monorepo

| Servicio        | Descripción                                        |
|------------------|----------------------------------------------------|
| `auth`          | Autenticación y validación de usuarios             |
| `user`          | Gestión de usuarios y relaciones follow/unfollow   |
| `tweet`         | Gestión de tweets y timelines con caching          |
| `gateway`       | Orquestador de llamadas entre microservicios       |

---

## 🏛️ Arquitectura

- **Hexagonal (Ports & Adapters):** Cada servicio está desacoplado y preparado para ser testeable.
- **Kafka:** Bus de eventos para comunicación asíncrona.
- **PostgreSQL:** Persistencia por servicio (base de datos independiente).
- **Redis:** Cache de timelines en `tweet` para mejorar tiempos de respuesta.
- **Docker Compose:** Entorno de desarrollo local completo (👉 **RECOMENDADO: levantar todo con Docker Compose**).
- **Go 1.23:** Lenguaje principal para todos los servicios.
- **Ent ORM:** Toda la persistencia usa Ent, lo que da un modelo de datos fuertemente tipado y migraciones seguras. 💥

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

### 1️⃣ Prerrequisitos
- Docker & Docker Compose (👉 **Imprescindible, todo el setup local fue pensado para Compose**)
- Go 1.23

### 2️⃣ Clonar el repositorio
```bash
git clone https://github.com/tuusuario/twittah.git
cd twittah
```

### 3️⃣ Levantar el entorno completo
```bash
docker-compose up --build
```
👉 **Nota:** El `docker-compose.yml` fue ajustado para manejar dependencias como Kafka, PostgreSQL y Redis de forma automática.


### 4️⃣ Acceder a Postman Collections

```bash
cd Twittah_Gateway.postman_collection.json
```

---

## 🧪 Testing

### Estado actual
✅ Tests unitarios implementados en `gateway`.  
⚠️ El resto de los servicios no tiene cobertura de tests aún (anotado como prioridad para próximas iteraciones).

### Ejecutar los tests disponibles
```bash
cd services/gateway
go test -cover -count=1 ./...
```

---

## 📝 Logging

El proyecto usa **logging estructurado con Uber Zap**.  
El package `common.Logger` inicializa y expone un singleton para que todos los servicios utilicen la misma configuración de logging.

---

## 📦 Directorio principal

```
services/
│
├── auth/       # Microservicio de autenticación
├── user/       # Microservicio de usuarios
├── tweet/      # Microservicio de tweets
├── gateway/    # Orquestador API
│
common/         # Utilidades compartidas (logger, config, etc.)
docker-compose.yml
```

---

## 🛠️ Tecnologías clave

| Componente        | Tecnología                 |
|--------------------|----------------------------|
| Backend            | Go 1.23                    |
| Comunicación       | Kafka (event-driven)       |
| Persistencia       | PostgreSQL + Ent ORM 💥     |
| Cache              | Redis                      |
| Contenerización    | Docker & Docker Compose    |
| Configuración      | Viper                      |

---

## 📦 Flujo de datos simplificado

1. Usuario envía request a `Gateway`.
2. `Gateway` valida y orquesta peticiones a otros servicios.
3. Eventos importantes son enviados por Kafka.
4. Servicios reaccionan y actualizan su estado (DB o cache).

---

## 🤝 Contribución

1. Forkea el repositorio.
2. Crea una rama (`feature/nueva-funcionalidad`).
3. Haz commit de tus cambios.
4. Abre un Pull Request.

---

## 📬 Contacto

**Nicolás Sánchez** – [LinkedIn](https://www.linkedin.com/in/mario-nahuel-nicolas-sanchez)  
📧 nico_dd@outlook.com.ar

---

## 🏁 Roadmap (Post-MVP)

### ✅ Calidad y pruebas
- Añadir cobertura de tests a todos los microservicios (unitarios >70%, e2e en flujos criticos).
- Mejorar los test existentes en Gateway para casos edge.
- Mejorar los tipos de respuesta http ante errores.
- Añadir tests de carga y estrés para validar escalabilidad.
- Implementar circuit breaker para servicios externos.

### 🔐 Seguridad
- Implementar sistema de Login/Signup. 
- Implementar sistema de autenticación con Bearer Tokens. 

### 🌐 API Gateway
  - Anteponer Kong a Gateway para un escalado rápido.
  - Rate limiting por usuario/IP.
  - Logging de acceso centralizado.
  - Validación de CORS.
  - Microservicios en redes privadas, Gateway expuesto como unica entrada.

### 📊 Observabilidad
- Integrar Prometheus + Grafana para métricas y dashboards.
- Añadir alertas (latencia alta, errores 5xx, caídas de Kafka/Redis).
- Logs estructurados enriquecidos con tracing IDs.

### 🚀 Escalabilidad y optimización
- Implementar fan-out on write total en TweetService para distribuir tweets a timelines cacheados.
- Cachear timelines paginados en Redis (`timeline:<user_id>:page:<n>`).
- Usar prefetch y pipelining en Redis para reducir latencia en cache-miss.
- Compresión de payloads Kafka (Snappy o LZ4).
- Considerar sharding en PostgreSQL para TweetService si la carga lo requiere.

### ☁️ Infraestructura y despliegue
- Despliegue en AWS ECS o EKS con Helm charts.
- Infraestructura como código (Terraform).
- Añadir CI/CD completo con GitHub Actions o GitLab CI (build, lint, test, deploy a staging/prod).

### 👨‍💻 Developer Experience
- Añadir pre-commit hooks (validación de formato, lint, tests).
- Documentación detallada por microservicio.

