//go:build encore_shell

package shellruntime

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/cockroachdb/errors"

	"encore.dev/appruntime/apisdk/api"
	"encore.dev/appruntime/shared/appconf"
	"encore.dev/appruntime/shared/logging"
	"encore.dev/shell/shellruntime/internal/daemon"
	"encore.dev/shell/shellruntime/internal/platform"
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
	ActiveTransport interface {
		http.RoundTripper
		TraceURL(traceID string) string
	}
}

func (a *ApiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = cloneRequest(req) // per RoundTripper contract

	// Remove Encore meta-data from the request
	for k := range req.Header {
		if strings.HasPrefix(k, "X-Encore-Meta-") {
			req.Header.Del(k)
		}
	}

	resp, err := a.ActiveTransport.RoundTrip(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Display the trace for the trace
	if traceID := resp.Header.Get("X-Encore-Trace-ID"); requestTracingEnabled && traceID != "" {
		if url := a.ActiveTransport.TraceURL(traceID); url != "" {
			logging.RootLogger.Trace().Str("call", req.URL.Path).Msg(url)
		}
	}

	return resp, nil
}

type localLoopBackTransport struct {
}

func (lb *localLoopBackTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the request to point to the app running within the local daemon
	localhost, err := daemon.GetListenAddressForApp(req.Context())
	if err != nil {
		if errors.Is(err, daemon.ErrAppNotRunning) {
			// don't wrap this, as we want the original message
			return nil, err
		}
		return nil, errors.Wrap(err, "could not get listen address for app on local system")
	}
	req.URL.Host = localhost
	req.Host = localhost

	// Make the request
	return http.DefaultClient.Do(req)
}

func (lb *localLoopBackTransport) TraceURL(traceID string) string {
	return fmt.Sprintf("http://localhost:9400/%s/envs/local/traces/%s", appconf.Runtime.AppID, traceID)
}

type encorePlatformProxyTransport struct {
	envName string
}

func (t *encorePlatformProxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Rewrite the path to go via the platform's API proxy
	req.URL.Scheme = "https"
	req.URL.Host = "api.encore.dev"
	req.URL.Path = fmt.Sprintf("/apps/%s/envs/%s/api/proxy%s", appconf.Runtime.AppID, t.envName, req.URL.Path)
	return platform.Do(req)
}

func (t *encorePlatformProxyTransport) TraceURL(traceID string) string {
	return fmt.Sprintf("https://app.encore.dev/%s/envs/%s/traces/%s", appconf.Runtime.AppID, t.envName, traceID)
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
