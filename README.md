# go-iam

**go-iam** is a lightweight, multi-tenant Identity and Access Management (IAM) server built in **Golang**. It provides robust authentication and fine-grained authorization for modern applications. With support for custom roles, third-party auth providers, and multi-client setups, `go-iam` gives you full control over access management in a scalable and modular way.

> ✅ Admin UI: [go-iam-ui](https://github.com/melvinodsa/go-iam-ui)  
> 🔐 Backend: [go-iam](https://github.com/melvinodsa/go-iam)

<img src=".github/go-iam.png" alt="drawing" width="400"/>

---

## ✨ Features

### 🔀 Multi-Tenancy

- Create and manage **Projects**.
- Each project is **isolated**, ensuring strict data boundaries.

### 🔐 Authentication Provider Integration

- Supports **Google OAuth** login.
- Multiple auth providers can be added (extensible).
- Supports **shared credentials** across clients for simplicity.

### 🧩 Client Management

- Manage multiple **clients** (apps) under the same project.
- Reuse the same auth credentials across clients.

### 🧱 Role-Based Access Control (RBAC)

- Define **resources** and group them into **roles**.
- Create **custom roles** and assign them to users.
- Resource-level access control.

### 🛠️ Admin UI

- Full-featured UI for managing users, roles, projects, and clients.
- Built for clarity and operational efficiency.
- [Go to Admin UI Repo](https://github.com/melvinodsa/go-iam-ui)

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

### Prerequisites

- Go 1.21+
- MongoDB
- Node.js + PNPM (for Admin UI)
- Google OAuth Credentials
- Redis (optional but recommended for performance)

---

### 🧪 Backend Setup

```bash
git clone https://github.com/melvinodsa/go-iam.git
cd go-iam
cp sample.env .env
go run main.go
```
