
# 🐦 Plataforma tipo Twitter – Gateway API

Bienvenido a la API Gateway de nuestra plataforma tipo Twitter. Aquí aprenderás cómo interactuar con los microservicios a través del **GatewayService**.

---

## 🚀 Levantar el proyecto
Asegúrate de tener Docker y Docker Compose instalados.

```bash
docker-compose up --build
```

Esto levantará:
- **Gateway** (http://localhost:8080) – No usa base de datos
- **UserService** (http://localhost:8082) – Usa Postgres para relaciones follow/unfollow
- **Kafka & Zookeeper** (eventos asíncronos)
- **Postgres** (solo para UserService)

La base de datos y la tabla `follows` se crean automáticamente con los scripts en `initdb/`.

---

---

## 📦 Endpoints disponibles

### 🔥 Autenticación
> **Nota:** En este MVP todos los usuarios son válidos. Envía el ID de usuario en el header `X-User-Id`.

---

### 📝 Publicar un Tweet
**POST** `/tweets`

```http
POST /tweets HTTP/1.1
Host: localhost:8080
Content-Type: application/json
X-User-Id: 123

{
  "content": "Hola mundo desde nuestra API!"
}
```

📌 **cURL:**
```bash
curl -X POST http://localhost:8080/tweets   -H "Content-Type: application/json"   -H "X-User-Id: 123"   -d '{"content": "Hola mundo desde nuestra API!"}'
```

---

### ➕ Seguir a un Usuario
**POST** `/follow/{target_user_id}`

```http
POST /follow/456 HTTP/1.1
Host: localhost:8080
X-User-Id: 123
```

📌 **cURL:**
```bash
curl -X POST http://localhost:8080/follow/456   -H "X-User-Id: 123"
```

---

### ➖ Dejar de seguir a un Usuario
**DELETE** `/follow/{target_user_id}`

```http
DELETE /follow/456 HTTP/1.1
Host: localhost:8080
X-User-Id: 123
```

📌 **cURL:**
```bash
curl -X DELETE http://localhost:8080/follow/456   -H "X-User-Id: 123"
```

---

### 📰 Ver Timeline
**GET** `/timeline`

```http
GET /timeline HTTP/1.1
Host: localhost:8080
X-User-Id: 123
```

📌 **cURL:**
```bash
curl http://localhost:8080/timeline   -H "X-User-Id: 123"
```

📦 **Respuesta ejemplo:**
```json
{
  "user_id": "123",
  "timeline": [
    {
      "tweet_id": "789",
      "author_id": "456",
      "content": "Hola mundo desde nuestra API!",
      "created_at": "2025-07-07T14:32:00Z"
    }
  ]
}
```

---

## 💡 Notas importantes
- El header `X-User-Id` es **obligatorio** para todas las requests.
- Todos los writes (publicar tweet, follow/unfollow) son **asíncronos** vía Kafka.
- **Solo UserService tiene base de datos** (Postgres). Gateway no almacena nada.
- Si borras los volúmenes de Docker, la base y tablas se recrean automáticamente gracias a los scripts `.sql` en `initdb/`.
- Puedes agregar más tablas para UserService agregando scripts SQL en ese directorio.

---

## 📖 Documentación adicional
Consulta el archivo [`business.txt`](./business.txt) para entender las decisiones arquitectónicas.
