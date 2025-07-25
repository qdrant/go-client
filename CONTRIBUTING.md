## Report Issues on GitHub [Issues](https://github.com/qdrant/go-client/issues)

We track public bugs and feature requests using GitHub issues. Please report by [opening a new issue](https://github.com/qdrant/go-client/issues/new).

**Effective Bug Reports** should include:

- A clear summary or background
- Steps to reproduce the issue
  - Be as specific as possible
  - Include sample code when possible
- What you expect to happen
- What happened
- Additional notes (e.g., why you think the issue occurs or solutions you’ve tried that didn’t work)

## Contributing Code

Follow these steps before submitting a pull request:

### Building the Project

```bash
go build ./...
```

This will download all dependencies and compile the project.

### Running Tests

All test files are in the `qdrant_test` directory and use [Testcontainers Go](https://golang.testcontainers.org/) for integration tests.

Run the following command to execute the test suites:

```bash
go test -v ./...
```

This command pulls a Qdrant Docker image to run integration tests. Ensure Docker is running.

### Formatting and Linting

Ensure your code is free from warnings and follows project standards.

The project uses [Gofmt](https://go.dev/blog/gofmt) for formatting and [golangci-lint](https://github.com/golangci/golangci-lint) for linting.

To format your code:

```bash
gofmt -s -w .
```

To lint your code:

```bash
golangci-lint run
```

### Preparing for a New Release

#### Pre-requisites

- [Go Installation](https://go.dev/doc/install)
- [Protobuf Compiler Installation](https://grpc.io/docs/protoc-installation/)

The client uses generated stubs from upstream Qdrant proto definitions, which are downloaded from [qdrant/qdrant](https://github.com/qdrant/qdrant/tree/master/lib/api/src/grpc/proto).

#### Steps

1. Download and generate the latest client stubs by running the following command from the project root:

```bash
BRANCH=dev sh internal/tools/sync_proto.sh
```

2. Update the `TestImage` value in [`qdrant_test/image_test.go`](https://github.com/qdrant/go-client/blob/master/qdrant_test/image_test.go) to `qdrant/qdrant:dev`.

3. Implement new Qdrant methods in [`points.go`](https://github.com/qdrant/go-client/blob/master/qdrant/points.go), [`collections.go`](https://github.com/qdrant/go-client/blob/master/qdrant/collections.go), or [`qdrant.go`](https://github.com/qdrant/go-client/blob/master/qdrant/qdrant.go) as needed and associated tests in [`qdrant_test/`](https://github.com/qdrant/go-client/tree/master/qdrant_test).

Since the API reference is published at <https://pkg.go.dev/github.com/qdrant/go-client>, the docstrings have to be appropriate.

4. If there are any new `oneOf` properties in the proto definitions, add helper constructors to [`oneof_factory.go`](https://github.com/qdrant/go-client/blob/master/qdrant/oneof_factory.go) following the existing patterns.

5. Run the linter, formatter and tests as per the instructions above.

6. Submit your pull request and get those approvals.

### Releasing a New Version

Once the new Qdrant version is live:

1. Run the following command:

```bash
BRANCH=master sh internal/tools/sync_proto.sh
```

2. Update the `TestImage` value in `qdrant_test/image_test.go` to `qdrant/qdrant:vNEW_VERSION`.

3. Merge the pull request.

4. Push a new Git tag to publish the version:

```bash
git tag v1.11.0
git push --tags
```

5. Optionally, do a release at <https://github.com/qdrant/go-client/releases> from the tag with notes.
