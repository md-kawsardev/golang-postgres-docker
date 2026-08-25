# Secure Notes API

A production-oriented REST API built with Go, PostgreSQL, JWT authentication, and Docker.

This project demonstrates how to build a secure backend application using a layered architecture, PostgreSQL data modeling, JWT-based authentication, ownership-based authorization, request validation, and containerized development.

The project was built as a practical backend engineering project with a focus on clean structure, security fundamentals, database design, and maintainability.

---

## Project Overview

The Secure Notes API allows authenticated users to create, view, and update their own private notes.

Each note belongs to a specific authenticated user, and the API enforces ownership at the database query level.

The project demonstrates a complete backend flow:

Client
↓
HTTP Request
↓
Gin Router
↓
JWT Authentication Middleware
↓
Controller
↓
Service Layer
↓
Repository Layer
↓
PostgreSQL

---

## Key Features

- User registration
- User login
- Secure password hashing
- JWT-based authentication
- JWT signature verification
- JWT expiration validation
- Protected API routes
- User identity extracted from JWT claims
- Create personal notes
- Retrieve authenticated user's notes
- Update notes with ownership verification
- PostgreSQL foreign key relationships
- Database constraints and indexes
- Request validation
- SQL injection prevention through parameterized queries
- Layered backend architecture
- Dockerized PostgreSQL environment
- Dockerized Go application
- Docker Compose development environment

---

## Tech Stack

### Backend

- Go
- Gin Web Framework
- JWT (`github.com/golang-jwt/jwt/v5`)
- pgx PostgreSQL driver
- pgx connection pool

### Database

- PostgreSQL
- SQL
- Foreign Keys
- Primary Keys
- Unique Constraints
- Indexes
- PostgreSQL Sequences

### DevOps / Infrastructure

- Docker
- Docker Compose
- Alpine Linux
- Multi-stage Docker builds

### Development Tools

- PowerShell
- Git
- GitHub
- REST API testing with PowerShell `Invoke-RestMethod`

---

## Architecture

The application follows a layered architecture:

```text
                 Client
                   |
                   v
            Gin HTTP Router
                   |
                   v
        Authentication Middleware
                   |
                   v
               Controller
                   |
                   v
                Service
                   |
                   v
               Repository
                   |
                   v
               PostgreSQL

```

### Responsibilities

#### Router

* Defining API routes
* Separating public and protected endpoints
* Applying authentication middleware

#### Middleware

* Reading the Authorization header
* Validating Bearer tokens
* Verifying JWT signatures
* Validating token expiration
* Extracting the authenticated user ID
* Rejecting invalid authentication attempts

#### Controller

* Receiving HTTP requests
* Parsing JSON request bodies
* Reading authenticated user identity
* Returning HTTP responses
* Mapping application errors to HTTP responses

#### Service

* Business logic
* Input validation
* Creating domain models
* Enforcing application-level rules

#### Repository

* PostgreSQL queries
* Database interaction
* Creating notes
* Finding notes by user
* Updating notes
* Keeping SQL/database logic separate from business logic

---

## Authentication

Authentication is implemented using JSON Web Tokens (JWT).

After successful login, the API generates a signed JWT containing the authenticated user's ID and token timestamps.

Example claims:

```json
{
  "user_id": 5,
  "iat": 1787590695,
  "exp": 1787677095
}

```

The token is signed using HMAC SHA-256.

### JWT Verification

Protected endpoints use authentication middleware. The middleware:

1. Reads the Authorization header
2. Validates the `Bearer <token>` format
3. Ensures the signing algorithm is HS256
4. Verifies the JWT signature using the configured JWT secret
5. Validates the token expiration
6. Extracts `user_id`
7. Stores the authenticated user ID in the Gin context

Example header:

```http
Authorization: Bearer <JWT_TOKEN>

```

Invalid or modified tokens are rejected.

---

## Authorization and Ownership

Authentication answers: **Who are you?**

Authorization answers: **Are you allowed to access this resource?**

This project implements ownership-based authorization. For example, when updating a note, the repository does not simply search by note ID. It uses both parameters:

```sql
WHERE id = $3 AND user_id = $4

```

This means a user can only update a note if:

* `note.id == requested_note_id` AND
* `note.user_id == authenticated_user_id`

The authenticated user ID comes from the verified JWT, not from the request body. This prevents a client from submitting another user's ID and attempting to modify their resources.

---

## Database Design

The project uses PostgreSQL with two main tables: `users` and `notes`.

### Users Table

```text
users
├── id (PK)
├── name
├── email (UNIQUE)
├── password_hash
├── created_at
└── updated_at

```

### Notes Table

```text
notes
├── id (PK)
├── title
├── content
├── user_id (FK -> users.id)
├── created_at
└── updated_at

```

The relationship uses:

```sql
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

```

This ensures that notes always belong to valid users.

---

## SQL Injection Prevention

Database queries use parameterized SQL through `pgx`.

Example:

```go
query := `
    SELECT id, name, email, password_hash, created_at, updated_at
    FROM users
    WHERE email = $1
`

err := r.db.QueryRow(
    ctx,
    query,
    email,
).Scan(...)

```

User input is passed separately as a query parameter rather than being concatenated into SQL strings.

---

## API Endpoints

### Authentication

| Method | Endpoint | Authentication |
| --- | --- | --- |
| `POST` | `/api/auth/register` | Public |
| `POST` | `/api/auth/login` | Public |

### Protected Routes

| Method | Endpoint | Authentication |
| --- | --- | --- |
| `GET` | `/api/test-protected` | JWT |
| `POST` | `/api/notes` | JWT |
| `GET` | `/api/notes` | JWT |
| `PUT` | `/api/notes/:id` | JWT |

---

## Example API Requests

### Register

```http
POST /api/auth/register
Content-Type: application/json

{
  "name": "Charlie",
  "email": "charlie@example.com",
  "password": "secret123"
}

```

### Login

```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "charlie@example.com",
  "password": "secret123"
}

```

Example response:

```json
{
  "message": "login successful",
  "token": "<JWT_TOKEN>"
}

```

### Create Note

```http
POST /api/notes
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "My First Secure Note",
  "content": "Learning PostgreSQL, JWT and SQL injection prevention."
}

```

### Get My Notes

```http
GET /api/notes
Authorization: Bearer <JWT_TOKEN>

```

### Update Note

```http
PUT /api/notes/2
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "title": "Updated Ownership Test",
  "content": "This note was updated securely."
}

```

---

## Project Structure

```text
golang-postgres-docker/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── config/
│   └── config.go
│
├── controllers/
│   ├── auth_controller.go
│   └── note_controller.go
│
├── database/
│   └── database.go
│
├── middleware/
│   └── auth_middleware.go
│
├── models/
│   ├── user.go
│   └── note.go
│
├── repository/
│   ├── user_repository.go
│   └── note_repository.go
│
├── routes/
│   └── routes.go
│
├── services/
│   ├── auth_service.go
│   └── note_service.go
│
├── utils/
│   └── jwt.go
│
├── migrations/
│   └── 001_schema.sql
│
├── Dockerfile
├── docker-compose.yml
├── .dockerignore
├── go.mod
├── go.sum
└── README.md

```

---

## Docker & Deployment

### Containerization

The project uses Docker to provide a consistent development environment. PostgreSQL runs inside a Docker container while exposing the port externally (`localhost:5433` maps to internal `5432`). A named volume is used for data persistence.

### Multi-Stage Docker Build

The Go application utilizes multi-stage builds:

* **Build Stage:** `golang:1.25-alpine` - Compiles the Go binary.
* **Runtime Stage:** `alpine:3.22` - Contains only the compiled binary for a minimal production image.

---

## Database Migrations

The database schema is defined in `migrations/001_schema.sql`. It contains:

* Tables
* Sequences
* Primary & Unique Keys
* Foreign Key Constraints
* Indexes & Defaults

---

## Environment Configuration

Sensitive configurations should be provided via environment variables (e.g., in a `.env` file):

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5433/note_db
JWT_SECRET=your-secret-key

```

---

## Running Locally

### Prerequisites

* Go installed
* Docker & Docker Compose installed
* Git

### Quickstart

1. **Clone the repository:**
```bash
git clone <repository-url>
cd golang-postgres-docker

```


2. **Install dependencies:**
```bash
go mod download

```


3. **Start PostgreSQL container:**
```bash
docker compose up -d

```


4. **Run the API server:**
```bash
go run ./cmd/server

```


The API will start at `http://localhost:8080`.

---

## Security Practices Demonstrated

* **Password Hashing:** Passwords are never stored in plaintext.
* **JWT Verification:** Protected endpoints strictly check signatures, algorithms, and expiration.
* **Ownership Verification:** Enforced at the repository layer using authenticated identity.
* **SQL Injection Defense:** Strict usage of parameterized queries.
* **Input Validation:** Mandatory fields validation before persistent operations.

---

## What This Project Demonstrates

* Go backend development & REST API design (Gin framework)
* PostgreSQL data modeling & relationship constraints
* Layered architecture (Controller, Service, Repository design patterns)
* JWT-based AuthN/AuthZ implementation
* Containerization with Docker & Multi-stage builds

---

## Future Improvements

* [ ] Delete & Get Single Note endpoints
* [ ] Refresh Tokens implementation
* [ ] Structured logging & Rate Limiting
* [ ] Automated Unit & Integration Tests
* [ ] OpenAPI / Swagger documentation
* [ ] CI/CD pipeline integration via GitHub Actions

---

## Author

Built as a practical backend engineering project to demonstrate experience with Go, PostgreSQL, REST API development, authentication, authorization, database security, and Docker.