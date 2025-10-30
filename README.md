# 🏧️ Multi-Tenant Go REST + gRPC API

![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14%2B-336791?logo=postgresql)
![Build](https://img.shields.io/badge/Build-Passing-brightgreen)

A **multi-tenant application** built with **Go (Golang)**, featuring both **REST** and **gRPC** interfaces.  
This project demonstrates a clean architecture using **GORM**, **PostgreSQL**, and **Google Wire** for dependency injection — with support for **multi-tenancy using schemas** (one master DB managing multiple client DBs).

---

## 🚀 Features

- ✅ REST API using [Fiber](https://github.com/gofiber/fiber)
- ✅ gRPC API with [grpc-go](https://github.com/grpc/grpc-go)
- ✅ Multi-Tenancy (Schema-based) — one schema per tenant
- ✅ GORM ORM with PostgreSQL
- ✅ Google Wire for Dependency Injection
- ✅ Environment-based configuration (dev, prod)
- ✅ Docker support
- ✅ Graceful shutdown & modular structure

---

## 🧩 Project Structure

```
go-multitenant/
├── cmd/
│ └── main.go
├── proto/
│ └── user.proto
├── internal/
│ ├── app/
│ │ ├── wire.go // wire injectors
│ │ └── wire_gen.go (generated)
│ ├── db/
│ │ ├── master_db.go
│ │ ├── db_provider.go
│ │ └── tx.go // transaction helper
│ ├── handler/
│ │ ├── rest.go
│ │ └── grpc_server.go
│ ├── service/
│ │ └── user_service.go
│ ├── repo/
│ │ └── user_repo.go
│ └── model/
│ ├── client.go
│ └── user.go
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Setup Guide

### 🧪 Prerequisites

- **Go** ≥ 1.22  
- **PostgreSQL** ≥ 14  
- **protoc** (Protocol Buffers compiler)  
- **Google Wire**  

Install dependencies:

```bash
brew install go protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/google/wire/cmd/wire@latest
```

---

### 🗄️ Database Setup

Start PostgreSQL:
```bash
brew services start postgresql
```

Login to Postgres:
```bash
psql -U admin
```

Create master database:
```sql
CREATE DATABASE masterdb;
```

Create sample tenant schemas:
```sql
CREATE SCHEMA tenant1;
CREATE SCHEMA tenant2;
```

---

### ⚙️ Environment Variables

Create a `.env` file in the project root:

```env
# Master Database (Tenant Registry)
MASTER_DSN=host=localhost user=postgres password=postgres dbname=masterdb port=5432 sslmode=disable

# Server Config
REST_PORT=8080
GRPC_PORT=50051
ENV=development
```

---

## 🧠 Dependency Injection (Wire)

Generate dependencies:
```bash
cd cmd
wire
```

---

## 🌐 Run REST API

```bash
go run cmd/main.go
```

### Example REST Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| GET | `/api/v1/books` | List books (tenant-aware) |
| POST | `/api/v1/books` | Create new book |
| GET | `/api/v1/books/:id` | Get book by ID |
| PUT | `/api/v1/books/:id` | Update book |
| DELETE | `/api/v1/books/:id` | Delete book |

> 🧩 Tenant is resolved dynamically using request header `X-Tenant-ID`.

---

## ⚡ Run gRPC Server

```bash
go run cmd/main.go
```

### Example `.proto` Definition

```protobuf
syntax = "proto3";
package userpb;
option go_package = "./userpb";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse) {}
}
message GetUserRequest {
    string client_id = 1;
    int64 user_id = 2;
}
message GetUserResponse {
    int64 user_id = 1;
    string name = 2;
}
```

Generate gRPC files:
```bash
protoc --go_out=. --go-grpc_out=. ./user.proto
```

---

## 🥪 Run REST + gRPC Together

### Option 1 — Run Separately
##### but here: go run ./cmd
```bash
go run cmd/rest/main.go
go run cmd/grpc/main.go

```

### Option 2 — Unified Injector (Advanced)
```bash
go run cmd/main.go
```

*(This uses a single Wire injector returning both servers)*

---

## 🐳 Docker Support

### Build
```bash
docker build -t multi-tenant-go .
```

### Run
```bash
docker run -p 8080:8080 -p 50051:50051 --env-file .env multi-tenant-go
```

---

## 🧰 Makefile Commands

| Command | Description |
|----------|-------------|
| `make wire` | Generate DI using Wire |
| `make proto` | Compile `.proto` files |
| `make run-rest` | Run REST API |
| `make run-grpc` | Run gRPC API |
| `make docker` | Build Docker image |

---

## 📖 Notes

- Multi-tenancy uses **PostgreSQL schemas**, not separate databases.  
- The **master database** stores tenant info (name, schema, credentials).  
- Tenant schema is dynamically switched at runtime.  
- Works for both **REST** and **gRPC** calls.  

---

---
## Step-by-Step Guide: Multi-Tenant Go Application Setup (Wire + GORM + Fiber + gRPC)
#### This guide walks through creating a production-ready multi-tenant Go application using Google Wire, GORM, Fiber, and gRPC.

## Install Required Tools
```bash
go install github.com/google/wire/cmd/wire@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
brew install protobuf # or sudo apt install protobuf-compiler
```

```bash
protoc --version
wire --version
```

## Bootstrap the Project
```bash
    mkdir go-multitenant
    cd go-multitenant
    go mod init go-multitenant
```
## Add dependencies:
```bash
    go get github.com/gofiber/fiber/v2
    go get github.com/google/wire
    go get gorm.io/gorm
    go get gorm.io/driver/postgres
    go get google.golang.org/grpc
    go get google.golang.org/protobuf
```
## Create Folder Structure
```bash
mkdir -p cmd internal/{app,db,handler,repo,service,model,proto}
```
---

## 👨‍💻 Author

**Devendra Pratap**  
Backend Engineer — Go | Spring Boot | PostgreSQL | Kafka  

[![GitHub](https://img.shields.io/badge/GitHub-DevendraPratap-black?logo=github)](https://github.com/devendrapratap307/)  
[![LinkedIn](https://img.shields.io/badge/LinkedIn-DevendraPratap-blue?logo=linkedin)](https://linkedin.com/in/devendrapratap307/)  

---



### 🌟 Star this repo if you find it useful!