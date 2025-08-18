# Simple Qdrant Client Example

This is a simple example application that demonstrates how to create a Go application that could interact with a Qdrant vector database server.

## What it does

The application:
1. Attempts to connect to a Qdrant server (default: localhost:6334)
2. Simulates checking if a collection named "test_collection" exists
3. Simulates creating the collection if it doesn't exist
4. Verifies the collection creation process
5. Demonstrates proper error handling and logging

## Features

- **Mock Qdrant Client**: Uses a simplified mock client to demonstrate the concept
- **Environment Configuration**: Configurable Qdrant server address via `QDRANT_HOST` environment variable
- **Proper Logging**: Comprehensive logging for debugging and monitoring
- **Error Handling**: Graceful handling of connection failures and errors
- **Docker Support**: Multi-stage Docker build for production deployment

## Prerequisites

- Docker installed on your system
- Go 1.22+ (for local development)

## Quick Start

### 1. Build and Run with Docker

```bash
# Make the build script executable (if not already)
chmod +x build-and-run.sh

# Build and run
./build-and-run.sh
```

### 2. Or Build and Run Manually

```bash
# Build the Docker image
docker build -t qdrant-simple-example .

# Run the container
docker run --rm qdrant-simple-example
```

### 3. Local Development

```bash
# Install dependencies
go mod tidy

# Build the application
go build -o main .

# Run the application
./main
```

## Configuration

You can customize the Qdrant server address by setting the `QDRANT_HOST` environment variable:

```bash
# Docker
docker run --rm -e QDRANT_HOST=your-qdrant-server:6334 qdrant-simple-example

# Local
export QDRANT_HOST=your-qdrant-server:6334
./main
```

## Expected Output

If successful, you should see output similar to:
```
2025/08/18 19:19:47 Connecting to Qdrant server at: localhost:6334
2025/08/18 19:19:47 Successfully connected to Qdrant server
2025/08/18 19:19:47 Checking if collection exists: test_collection
2025/08/18 19:19:47 Attempting to create collection: test_collection
2025/08/18 19:19:48 Successfully created collection: test_collection
2025/08/18 19:19:48 Successfully created collection 'test_collection'
2025/08/18 19:19:48 Checking if collection exists: test_collection
2025/08/18 19:19:48 Collection 'test_collection' creation completed
2025/08/18 19:19:48 Example completed successfully!
```

## Architecture

- **MockQdrantClient**: A simplified client that demonstrates the interface pattern
- **gRPC Connection**: Uses Google's gRPC library for network communication
- **Context Management**: Proper timeout handling with context
- **Error Handling**: Comprehensive error handling and logging

## Docker Details

The Dockerfile uses a multi-stage build approach:
1. **Builder Stage**: Compiles the Go application
2. **Final Stage**: Creates a minimal Alpine Linux container with just the binary
3. **Security**: Runs as a non-root user for better security
4. **Size**: Optimized for minimal container size

## Files

- `main.go` - The main application code with mock Qdrant client
- `go.mod` - Go module definition with external dependencies
- `go.sum` - Dependency lock file
- `Dockerfile` - Multi-stage Docker build configuration
- `build-and-run.sh` - Convenience script for building and running
- `README.md` - This documentation file

## Next Steps

To extend this example with a real Qdrant client:
1. Replace the mock client with the actual `github.com/qdrant/go-client/qdrant` package
2. Update the `go.mod` to include the real Qdrant client dependency
3. Modify the client calls to use the actual Qdrant API methods
4. Test with a real Qdrant server instance
