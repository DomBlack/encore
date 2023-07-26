//go:build encore_shell

package shellruntime

import (
	"net/http"

	"github.com/cockroachdb/errors"

	"encore.dev/appruntime/apisdk/api"
)

func init() {
	// Override the API framework's HTTP client with our own
	// which will route the request to the correct environment
	api.Singleton.SetHTTPClient(&http.Client{
		Transport: shellApiTransportSingleton,
	})
}

var shellApiTransportSingleton = &ApiTransport{
	ActiveTransport: &localLoopBackTransport{},
}

type ApiTransport struct {
	ActiveTransport http.RoundTripper
}

func (a *ApiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return a.ActiveTransport.RoundTrip(req)
}

type localLoopBackTransport struct {
}

func (lb *localLoopBackTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("not implemented")
}
