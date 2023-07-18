//go:build encore_shell

package platform

import (
	"context"

	"github.com/cockroachdb/errors"
)

// Query will execute a GraphQL query against the platform.
func Query(ctx context.Context, query string, variables map[string]any, resp any) error {
	if defaultClient == nil {
		return errors.New("no platform client configured")
	}

	return defaultClient.Query(ctx, query, variables, resp)
}
