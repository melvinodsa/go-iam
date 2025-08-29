# go-iam

[![codecov](https://codecov.io/github/melvinodsa/go-iam/graph/badge.svg?token=TWJXNBHTQL)](https://codecov.io/github/melvinodsa/go-iam)
[![Go Report Card](https://goreportcard.com/badge/github.com/melvinodsa/go-iam)](https://goreportcard.com/report/github.com/melvinodsa/go-iam)
[![GoDoc](https://godoc.org/github.com/melvinodsa/go-iam?status.svg)](https://godoc.org/github.com/melvinodsa/go-iam)

**go-iam** is a lightweight, multi-tenant Identity and Access Management (IAM) server built in **Golang**. It provides robust authentication and fine-grained authorization for modern applications. With support for custom roles, third-party auth providers, and multi-client setups, `go-iam` gives you full control over access management in a scalable and modular way.

> ✅ Admin UI: [go-iam-ui](https://github.com/melvinodsa/go-iam-ui)  
> 🐳 Docker Setup: [go-iam-docker](https://github.com/melvinodsa/go-iam-docker)  
> 🔐 Backend: [go-iam](https://github.com/melvinodsa/go-iam)  
> 📦 SDK: [go-iam-sdk](https://github.com/melvinodsa/go-iam-sdk)  
> 🚀 Examples: [go-iam-examples](https://github.com/melvinodsa/go-iam-examples)

<img src=".github/go-iam.png" alt="go-iam overview" width="400"/>

---

## ✨ Features

### 🔀 Multi-Tenancy

- Create and manage **Projects**
- Strict **isolation** of data between tenants

### 🔐 Authentication Provider Integration

- Google OAuth login support
- Easily extendable to add more providers
- **Shared credentials** support across multiple clients

### 🧩 Client Management

- Multiple apps (clients) per project
- Avoid duplicate OAuth credentials

### 🧱 Role-Based Access Control (RBAC)

- Define resources and group them into roles
- Create custom roles and assign to users
- Granular access control for different actions/resources

### 🛠️ Admin UI

- React-based Admin interface for managing:
  - Projects
  - Users
  - Roles
  - Resources
  - Clients

---

## 🧰 Tech Stack

| Component     | Tech                |
| ------------- | ------------------- |
| Backend       | Golang              |
| Database      | MongoDB             |
| Caching (opt) | Redis               |
| Frontend      | React + Vite (PNPM) |

---

## 🚀 Getting Started

### Option 1: 🔧 Manual Setup (Development)

#### Prerequisites

- Go 1.21+
- MongoDB
- Redis (optional, recommended)
- Google OAuth Credentials

#### Run the Backend

```bash
git clone https://github.com/melvinodsa/go-iam.git
cd go-iam
cp sample.env .env
go run main.go
```

### Option 2: 🐳 Docker-Based Local Setup (Recommended for Testing)

Use the official go-iam-docker repo to spin up everything with Docker Compose, including:

- MongoDB
- Redis
- go-iam (backend)
- go-iam-ui (admin frontend)

#### Steps

```bash
git clone https://github.com/melvinodsa/go-iam-docker.git
cd go-iam-docker
cp sample.env .env
docker compose up -d
```

#### Access

- Admin UI: [http://localhost:4173](http://localhost:4173)
- API: [http://localhost:3000](http://localhost:3000)
- API Docs: [http://localhost:3000/docs](http://localhost:3000/docs)

## 📦 Environment Variables

Some important environment variables used in `.env`:

| Variable                                       | Description                                                           |
| ---------------------------------------------- | --------------------------------------------------------------------- |
| `LOGGER_LEVEL`                                 | Logger level `1 - Debug` (refer., `https://docs.gofiber.io/api/log/`) |
| `DB_HOST`                                      | MongoDB URI (e.g., `mongodb://user:pass@host/db`)                     |
| `JWT_SECRET`                                   | Secret key used for generating and verifying JWT tokens               |
| `REDIS_HOST`, `REDIS_PASSWORD`, `ENABLE_REDIS` | Redis host address and toggle to enable Redis caching                 |
| `ENCRYPTER_KEY`                                | Optional symmetric key for encrypting sensitive fields - change this  |
| `AUTH_PROVIDER_REFETCH_INTERVAL_IN_MINUTES`    | Interval in minutes to refetch and sync third-party auth providers    |
| `TOKEN_CACHE_TTL_IN_MINUTES`                   | Interval for which the authentication token should be valid           |

## License

- Community Edition: [Apache 2.0](./LICENSE) (Open Source, free to use)
