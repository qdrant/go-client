# Simple Qdrant Client Example

This is a simple example application that demonstrates how to create a Go application that could interact with a Qdrant vector database server.

## Prerequisites

- Docker installed on your system

## Quick Start

### 1. Start Qdrant

First, make sure Qdrant is running. You can start it with Docker:

```bash
docker run --rm -it -p 6334:6334 -p 6333:6333 qdrant/qdrant
```

### 2. Building and Running

```bash
bash -x build-and-run.sh
```
