FROM golang:1.22-alpine

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY main.go ./

# Build and run the application
CMD ["go", "run", "main.go"]
