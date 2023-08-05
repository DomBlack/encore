// Package daemon is used to talk to the local daemon on the host computer
package daemon

import (
	"context"
	goErrs "errors"
	"strings"
	"sync"

	"github.com/cockroachdb/errors"

	"encore.dev/appruntime/shared/appconf"
)

//go:generate protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./mini_daemon.proto

var (
	listenAddressMu     sync.Mutex
	listenAddressCached string
)

var (
	ErrAppNotRunning = goErrs.New("app not running, please start it using `encore run`")
)

// GetListenAddressForApp returns the listen address for the currently running instance of the app
func GetListenAddressForApp(ctx context.Context) (string, error) {
	listenAddressMu.Lock()
	defer listenAddressMu.Unlock()

	if listenAddressCached != "" {
		return listenAddressCached, nil
	}

	client, err := dial(ctx)
	if err != nil {
		return "", err
	}

	resp, err := client.ListenAddressForApp(ctx, &ListenAddressForAppRequest{
		AppId: appconf.Runtime.AppID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "app not running") {
			return "", errors.WithStack(ErrAppNotRunning)
		}
		return "", errors.Wrap(err, "could not get listen address for app on local system")
	}
	if resp.ListenAddress == "" {
		return "", errors.New("no listen address returned from Encore daemon")
	}
	listenAddressCached = resp.ListenAddress

	return listenAddressCached, nil
}
