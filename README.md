# Go + Fiber Hello World Project

This project creates a simple web server using Go and the Fiber framework that returns a JSON response.

## Prerequisites

Make sure you have Go installed on your system. You can download it from: https://go.dev/dl/

## Installation

1. Install Go from https://go.dev/dl/ (minimum version 1.17)
2. Navigate to this project directory
3. Run the following commands:

```bash
go mod tidy
go run main.go
```

## Usage

Once the server is running, you can access it at:

```
http://localhost:3000/
```

Expected response:
```json
{"message": "hello world"}
```

## Project Structure

- `main.go` - Main application file with Fiber server setup
- `go.mod` - Go module file with dependencies
- `README.md` - This file

## API Endpoints

- `GET /` - Returns JSON with hello world message