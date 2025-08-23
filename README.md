# GoLang Full-Stack gRPC Project

This project is a full-stack web application built entirely in Go, featuring a microservices architecture.

## Tech Stack

- **Backend**: Go, gRPC, gRPC-Gateway
- **Frontend**: Go, Echo Framework, HTML Templates, JavaScript (`fetch`)
- **Database**: PostgreSQL
- **ORM**: Bun ORM

## Prerequisites

1.  **Go**: Version 1.21 or later.
2.  **PostgreSQL**: A running instance of PostgreSQL.
3.  **Protocol Buffers Compiler (`protoc`)**:
    - `protoc`
    - `protoc-gen-go`
    - `protoc-gen-go-grpc`
    - `protoc-gen-grpc-gateway`

    Install the Go plugins:
    ```sh
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    go install [github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest](https://github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest)
    ```

## Project Setup

1.  **Clone the repository** (or create the files as provided).

2.  **Create a database** in PostgreSQL. For example, `fullstack_db`.

3.  **Set Environment Variables**:
    Create a `.env` file in the `backend` directory with your database credentials:
    ```env
    # backend/.env
    DB_DSN="postgres://YOUR_USER:YOUR_PASSWORD@localhost:5432/fullstack_db?sslmode=disable"
    ```

4.  **Install Go dependencies**:
    ```sh
    go mod tidy
    ```

## How to Run

### 1. Generate Protobuf Code

Navigate to the `backend` directory and run the `proto` command. This only needs to be done once or when you change the `.proto` files.

```sh
cd backend
make proto