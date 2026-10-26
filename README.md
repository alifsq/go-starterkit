# Go Backend Starter Kit V.1.0.0

A production-ready and modular Go backend starter kit built with clean architecture principles. Designed for scalability, maintainability, and quick project setup.

## 🚀 Features

- **Framework & ORM**: Built using [Gin](https://github.com/gin-gonic/gin) and [GORM Gen](https://github.com/go-gorm/gen) for type-safe database queries.
- **Dependency Injection**: Integrated with [Google Wire](https://github.com/google/wire) for compile-time dependency injection.
- **Database Migrations**: Structured and version-controlled database schema migrations.
- **Logging**: High-performance structured logging powered by [Uber Zap](https://github.com/uber-go/zap).
- **Authentication**: JWT-based authentication featuring access token, refresh token workflows, and session revocation.
- **Email Delivery**: Automated email verification system with time-limited tokens sent via SMTP.
- **Configuration**: Centralized configuration management using [Viper](https://github.com/spf13/viper).
- **Security & Access Control**:
  - **Role-Based Access Control (RBAC)** powered by Casbin for fine-grained permissions
  - **Rate Limiting** backed by Redis to prevent API abuse and brute-force attacks HttpOnly Cookie protection for sensitive tokens

## 🛠️ Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **Database**: MySQL / PostgreSQL
- **ORM**: GORM & GORM Gen
- **Logging**: Uber Zap
- **Authentication**: JWT (JSON Web Tokens)
- **Dependency Injection**: Google Wire
- **Configuration**: Viper
- **Email Service**: SMTP
- **Role Based Access Control(RBAC)**: Casbin

## 📦 Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (v1.27)
- [MySQL](https://www.mysql.com/) database server

### Installation

1. **Clone the repository**

   ```bash
   git clone [https://github.com/](https://github.com/)<your-username>/<your-repo-name>.git
   cd <your-repo-name>

   ```

2. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```
3. **Install dependencies**
   ```bash
   go mod download
   ```
4. **Run Database Migrations**
   ```bash
   make migrate-up
   ```
5. **Run the Application**
   ```bash
   make run
   ```

## 🛠️ Useful Commands (Makefile)

## 🛠️ Useful Commands (Makefile)

This project includes a `Makefile` to simplify common commands. Run the following commands in your terminal:

- **`make dev`** - Run the application directly in development mode.
- **`make build`** - Compile the Go code into an executable binary inside the `bin/` directory.
- **`make run`** - Compile and execute the binary file.
- **`make wire`** - Regenerate dependency injection code (`wire_gen.go`).
- **`make migrate-up`** - Run database schema migrations.
- **`make migrate-down`** - Rollback the last database migration.
- **`make test`** - Run all unit tests within the repository.
- **`make clean`** - Remove build binaries and temporary test artifacts.

---

## 💻 Built With

<p align="center">
  <em>Built with:</em><br><br>
  <img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Gin-008080?style=for-the-badge&logo=gin&logoColor=white" alt="Gin" />
  <img src="https://img.shields.io/badge/MySQL-4479A1?style=for-the-badge&logo=mysql&logoColor=white" alt="MySQL" />
  <img src="https://img.shields.io/badge/GORM-00599C?style=for-the-badge&logo=go&logoColor=white" alt="GORM" />
  <img src="https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens" alt="JWT" />
  <img src="https://img.shields.io/badge/Uber_Zap-000000?style=for-the-badge&logo=uber&logoColor=white" alt="Uber Zap" />
  <img src="https://img.shields.io/badge/Google_Wire-4285F4?style=for-the-badge&logo=google&logoColor=white" alt="Google Wire" />
  <img src="https://img.shields.io/badge/Viper-5A29E4?style=for-the-badge&logo=go&logoColor=white" alt="Viper" />
  <img src="https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
  <img src="https://img.shields.io/badge/Casbin-FF4500?style=for-the-badge&logo=security&logoColor=white" alt="Casbin" />
</p>
