# Simple WebSocket Chat with Fiber

A basic real-time chat application built with Go and Fiber, using WebSockets for communication.

## Features

- Real-time messaging using WebSockets
- Username support
- System messages for user join/leave events
- Modern styled web interface
- Console client for terminal users
- Docker support for easy deployment

## Prerequisites

- Go 1.24 or higher
- Docker (optional, for containerized deployment)

## Project Structure

```
.
├── cmd
│   ├── client
│   │   └── main.go     // Client entry point
│   └── server
│       └── main.go     // Server entry point
├── internal
│   ├── client
│   │   └── client.go   // Client implementation
│   └── server
│       └── server.go   // Server implementation with Fiber
├── .env                // Environment variables
├── Dockerfile          // Server Docker configuration
├── Dockerfile.client   // Client Docker configuration
├── docker-compose.yml  // Docker Compose configuration
├── go.mod              // Go module file
└── go.sum              // Go dependencies checksums
```

## Configuration

The application uses environment variables for configuration which can be set in the `.env` file:

```env
# Server address
SERVER_ADDR=0.0.0.0:8080

# WebSocket URL for the client
WS_URL=ws://localhost:8080/ws
```

## Running Locally

1. Run the server:
   ```
   go run cmd/server/main.go
   ```

2. In another terminal, run the client:
   ```
   go run cmd/client/main.go
   ```

3. Enter your username when prompted and start chatting!

## Running with Docker

### Using Docker Compose (recommended)

```bash
docker-compose up --build
```

### Running the Server Only

```bash
docker build -t chat-server .
docker run -p 8080:8080 chat-server
```

### Running the Client

```bash
docker build -t chat-client -f Dockerfile.client .
docker run -it --network host chat-client
```

## WebSocket API

The server exposes a WebSocket endpoint at `/ws` that accepts the following query parameters:

- `username`: The user's name (optional)

## License

This project is open-source.