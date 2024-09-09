package qdrant

import (
	"fmt"
	"strings"
)

//nolint:revive // The linter says qdrant.QdrantError stutters, but it's an apt name.
type QdrantError struct {
	operationName string
	context       string
	err           error
}

// Error returns the error as string.
func (e *QdrantError) Error() string {
	if e.context == "" {
		return fmt.Sprintf("%s() failed: %v", e.operationName, e.err)
	}
	return fmt.Sprintf("%s() failed: %s: %v", e.operationName, e.context, e.err)
}

// Unwrap returns the inner error.
func (e *QdrantError) Unwrap() error {
	return e.err
}

func newQdrantErr(err error, operationName string, contexts ...string) *QdrantError {
	combinedContext := strings.Join(contexts, ": ")
	return &QdrantError{
		operationName: operationName,
		err:           err,
		context:       combinedContext,
	}
}
