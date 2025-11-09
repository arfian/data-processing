# data-processing
data-processing

A modern application API built with Go and asynchronous CSV processing

## Prerequisites

Before running the application, make sure you have the following installed:

- Go 1.21 or higher
- PostgreSQL
- Golang Migrate

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/arfian/data-processing.git
cd data-processing
```

### 2. Environment Setup

Create a `.env` file in the root directory, you can copy value .env.example

```bash
WORKER_COUNT=5
BATCH_SIZE=20
DATABASE_URL=
SERVER_PORT=8088
```

### 3. Go-migrate CLI
```sh
#mac
$ brew install golang-migrate

#linux
$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz
```

### 4. Running Migration
#### Migration Up
```sh
$ make migrateup
```

#### Migration Down
```sh
$ make migratedown
```

## How To Run
### Using Makefile
```sh
$ make run 
```

### Using Terminal / Cmd
```sh
$ go mod download
$ go run main.go 
```

### Check Unit Test
```sh
$ make test 
```

## API Documentation
The API documentation is available in Postman format. Import the following files into Postman:

- `postman/Data Process.postman_collection.json`

### Key Endpoints

1. CSV
   - POST `/api/v1/csv/process` - User registration

## Project Structure
```
.
├── config/           # Setup config register env application
├── csv/              # Folder File CSV
├── docs/             # Generate docs API swagger
├── internal/         # Code with internal system
│   ├── delivery/     # Bussiness domain
│   ├── ├── http/     # Logic handle API logic
│   ├── domain/       # List all interface domain
│   ├── repository/   # Logic query SQL
│   ├── usecase/      # Business logic implementation
└── migrations/       # Generate code migration database
└── pkg/              # code package helper logic
└── postman/          # Postman collection and environment
```

## Technologies
- [Golang](https://go.dev/)
- [Gorm](https://gorm.io/index.html)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swaggo](https://github.com/swaggo/swag)
- [Gin](https://gin-gonic.com/)
- PostgreSQL

## Accessing Swagger
```
localhost:8089/swagger/index.html
```