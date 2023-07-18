//go:build encore_shell

package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
)

// Query executes a GraphQL query against the Encore platform
func (c *Client) Query(ctx context.Context, query string, variables map[string]any, respObj any) error {
	// Create the request
	req := graphQLRequest{Query: query, Variables: variables}
	reqBytes, _ := json.Marshal(req)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", "/graphql", bytes.NewReader(reqBytes))
	if err != nil {
		return errors.Wrap(err, "failed to create GraphQL request")
	}
	httpReq.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := c.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "failed to execute GraphQL request")
	}
	defer func() { _ = resp.Body.Close() }()

	// Decode the response
	graphResponse := &graphQLResponse{}
	if err := json.NewDecoder(resp.Body).Decode(graphResponse); err != nil {
		return errors.Wrap(err, "failed to decode GraphQL response")
	} else if graphResponse.Errors != nil && len(*graphResponse.Errors) > 0 {
		return errors.WithStack(graphResponse.Errors)
	} else if respObj != nil {
		if err := json.Unmarshal(graphResponse.Data, respObj); err != nil {
			return errors.Wrap(err, "failed to decode GraphQL response data")
		}
	}
	return nil
}

type graphQLRequest struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables,omitempty"`
	OperationName string         `json:"operationName,omitempty"`
	Extensions    map[string]any `json:"extensions,omitempty"`
}

type graphQLResponse struct {
	Data       json.RawMessage   `json:"data"`
	Errors     *GraphQLErrorList `json:"errors,omitempty"`
	Extensions map[string]any    `json:"extensions,omitempty"`
}

// GraphQLError is an error returned by the GraphQL API.
type GraphQLError struct {
	Message    string                     `json:"message"`
	Path       []string                   `json:"path"`
	Extensions map[string]json.RawMessage `json:"extensions"`
}

func (e *GraphQLError) Error() string {
	return e.Message
}

// GraphQLErrorList is a list of GraphQLError objects.
type GraphQLErrorList []*GraphQLError

func (err GraphQLErrorList) Error() string {
	if len(err) == 0 {
		return "no errors"
	} else if len(err) == 1 {
		return err[0].Error()
	}
	return fmt.Sprintf("%s (and %d more errors)", err[0].Error(), len(err)-1)
}
