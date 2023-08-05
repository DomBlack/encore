//go:build encore_shell

package platform

import (
	"context"
	"net/http"

	"github.com/cockroachdb/errors"
)

// Query will execute a GraphQL query against the platform.
func Query(ctx context.Context, query string, variables map[string]any, resp any) error {
	if defaultClient == nil {
		return errors.New("no platform client configured")
	}

	return defaultClient.Query(ctx, query, variables, resp)
}

// Do will execute an HTTP request against the platform API
func Do(req *http.Request) (*http.Response, error) {
	if defaultClient == nil {
		return nil, errors.New("no platform client configured")
	}

	return defaultClient.RawDo(req)
}
