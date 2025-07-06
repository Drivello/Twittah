
# 🐦 Twittah – Backend Challenge Ualá

Este proyecto es una versión simplificada de una plataforma de microblogging tipo Twitter. Permite a los usuarios publicar tweets, seguir a otros usuarios y ver un timeline con los tweets de las personas que siguen.

> 🚀 Diseñado para escalar a millones de usuarios y optimizado para lecturas masivas.

---

## 📦 Features principales
- 📝 **Publicación de tweets** (máximo 280 caracteres)
- ➕ **Seguir y dejar de seguir usuarios**
- 📰 **Timeline optimizado** para lecturas en tiempo real
- ⚡ Arquitectura **Hexagonal (Ports & Adapters)** para mantener separación clara entre dominio y capa de infraestructura
- 📡 Comunicación entre microservicios mediante **Kafka**
- 🗄️ Persistencia con PostgreSQL y cache con Redis
- 🐳 Contenerización completa con Docker y orquestación con Docker Compose

---

## 📂 Arquitectura general
El sistema está compuesto por 4 microservicios desacoplados:

| Servicio         | Descripción                                  |
|-------------------|----------------------------------------------|
| GatewayService    | Orquestador y punto de entrada de la API REST|
| AuthService       | Valida usuarios en cada request              |
| UserService       | Maneja relaciones de seguimiento             |
| TweetService      | Maneja publicación y distribución de tweets  |

📝 Comunicación **asíncrona** en writes (Kafka).  
📰 Lecturas en tiempo real (Redis cache).  

---

## 🏗️ Diagrama de Arquitectura

```plaintext
                        ┌─────────────┐
                        │   Cliente   │
                        └──────▲──────┘
                               │ REST API
                      ┌────────┴────────┐
                      │ GatewayService  │
                      └────────┬────────┘
          ┌────────────────────┼────────────────────┐
          ▼                    ▼                    ▼
   ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
   │ AuthService  │     │ UserService  │     │ TweetService │
   └──────▲───────┘     └──────▲───────┘     └──────▲───────┘
          │ Kafka Events       │ Kafka Events       │ Kafka Events
          ▼                    ▼                    ▼
     PostgreSQL           PostgreSQL           PostgreSQL
                              │                      │
                              └──────────┬───────────┘
                                         ▼
                                      Redis
                            (Cache timelines por usuario)
```                           

---

## 🚀 Como utilizar esta APP

Detalles en [`HowToUseThisApp.md`](./HowToUseThisApp.md)

---

## ✅ Consideraciones
# Consideraciones y suposiciones detrás de la arquitectura de esta aplicación.

📖 Más detalles en [`business.txt`](./business.txt).

---

## 🏛️ Tecnologías
- **Lenguaje:** Golang
- **Persistencia:** PostgreSQL
- **Cache:** Redis
- **Mensajería:** Apache Kafka
- **Contenedores:** Docker + Docker Compose
- **Arquitectura:** Hexagonal (Ports & Adapters)

---

## 🧪 Testing
- Cobertura mínima del 50% en tests unitarios
- Tests de integración en GatewayService y consumidores Kafka

---

## 📬 Contacto

- Linkedin: https://www.linkedin.com/in/mario-nahuel-nicolas-sanchez/
- Email: Nico_DD@outlook.com.ar

