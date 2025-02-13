package contextual

import (
	"context"
	"fmt"

	"github.com/Carbonfrost/autogun/pkg/workspace"
)

type key string

const (
	workspaceKey key = "workspace"
)

// With will update the context with the given values.
func With(ctx context.Context, values ...any) context.Context {
	for _, v := range values {
		ctx = context.WithValue(ctx, keyFor(v), v)
	}
	return ctx
}

// Workspace gets the Workspace from the context otherwise panics
func Workspace(ctx context.Context) *workspace.Workspace {
	if res, ok := ctx.Value(workspaceKey).(*workspace.Workspace); ok {
		return res
	}
	panic(failedMust(workspaceKey))
}

func keyFor(a any) key {
	switch a.(type) {
	case *workspace.Workspace:
		return workspaceKey
	}
	panic(fmt.Errorf("unsupported context type %T", a))
}

func failedMust(k key) string {
	return fmt.Sprintf("expected %s value not present in context", k)
}
