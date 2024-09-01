package qdrant

import (
	"fmt"
	"strings"
)

//nolint:revive // The linter says qdrant.QdrantErr stutters, but it's an apt name.
type QdrantErr struct {
	operationName string
	context       string
	err           error
}

// Error returns the error as string
func (e *QdrantErr) Error() string {
	if e.context == "" {
		return fmt.Sprintf("%s() failed: %v", e.operationName, e.err)
	}
	return fmt.Sprintf("%s() failed: %s: %v", e.operationName, e.context, e.err)
}

// Unwrap returns the inner error
func (e *QdrantErr) Unwrap() error {
	return e.err
}

func newQdrantErr(err error, operationName string, contexts ...string) *QdrantErr {
	combinedContext := strings.Join(contexts, ": ")
	return &QdrantErr{
		operationName: operationName,
		err:           err,
		context:       combinedContext,
	}
}
