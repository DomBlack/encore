//go:build encore_shell

package platform

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"runtime"

	"github.com/cockroachdb/errors"

	"encore.dev/appruntime/shared/appconf"
	"encore.dev/shell/shellruntime/internal/platform/auth"
)

type Client struct {
	baseURL *url.URL
	client  *http.Client
}

var defaultClient *Client

// Init will initialize the default platform client.
func Init() error {
	authClient, err := auth.NewClient()
	if err != nil {
		return errors.Wrap(err, "failed to initialize Encore auth client")
	}

	baseURL, err := url.Parse(auth.APIBaseURL)
	if err != nil {
		return errors.Wrap(err, "failed to parse Encore API URL")
	}

	defaultClient = &Client{
		client:  authClient,
		baseURL: baseURL,
	}

	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.RawDo(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		defer func() { _ = resp.Body.Close() }()
		return nil, decodeErrorResponse(resp)
	}

	return resp, nil
}

func (c *Client) RawDo(req *http.Request) (*http.Response, error) {
	// Add a very limited amount of information for diagnostics
	req.Header.Set("User-Agent", "EncoreShell/"+appconf.Static.EncoreCompiler+"/"+appconf.Runtime.AppID)
	req.Header.Set("X-Encore-Version", appconf.Static.EncoreCompiler)
	req.Header.Set("X-Encore-GOOS", runtime.GOOS)
	req.Header.Set("X-Encore-GOARCH", runtime.GOARCH)

	req.URL = c.baseURL.ResolveReference(req.URL)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return resp, nil
}

func decodeErrorResponse(resp *http.Response) error {
	var respStruct struct {
		OK    bool
		Error Error
		Data  json.RawMessage
	}
	if err := json.NewDecoder(resp.Body).Decode(&respStruct); err != nil {
		return errors.Wrap(err, "decode response")
	}
	e := respStruct.Error
	e.HTTPCode = resp.StatusCode
	e.HTTPStatus = resp.Status
	return errors.WithStack(e)
}

type Error struct {
	HTTPStatus string `json:"-"`
	HTTPCode   int    `json:"-"`
	Code       string
	Detail     json.RawMessage
}

func (e Error) Error() string {
	if len(e.Detail) > 0 {
		return fmt.Sprintf("http %s: code=%s detail=%s", e.HTTPStatus, e.Code, e.Detail)
	}
	return fmt.Sprintf("http %s: code=%s", e.HTTPStatus, e.Code)
}
