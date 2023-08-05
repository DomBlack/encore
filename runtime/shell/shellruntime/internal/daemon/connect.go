package daemon

import (
	"context"
	"net"
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"

	"encore.dev/appruntime/shared/syncutil"
)

var (
	setupConn syncutil.Once
	conn      DaemonClient // Don't reference directly, use [dial].
)

// dail returns a connection to the Encore daemon.
func dial(ctx context.Context) (DaemonClient, error) {
	err := setupConn.Do(func() error {
		socketPath, err := daemonSockPath()
		if err != nil {
			return errors.Wrap(err, "could not determine daemon socket path")
		}

		dailer := func(ctx context.Context, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "unix", socketPath)
		}

		newConn, err := grpc.DialContext(ctx, "unix",
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithContextDialer(dailer),
		)
		if err != nil {
			return errors.Wrap(err, "could not dial daemon")
		}

		conn = NewDaemonClient(newConn)
		return nil
	})

	return conn, err
}

// daemonSockPath reports the path to the Encore daemon unix socket.
func daemonSockPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", errors.Wrap(err, "could not determine cache dir")
	}
	return filepath.Join(cacheDir, "encore", "encored.sock"), nil
}
