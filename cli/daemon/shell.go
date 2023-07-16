package daemon

import (
	"strings"

	"encr.dev/cli/daemon/run"
	"encr.dev/internal/optracker"
	daemonpb "encr.dev/proto/encore/daemon"
)

// Shell builds an interactive shell application for the given app.
func (s *Server) Shell(req *daemonpb.ShellRequest, stream daemonpb.Daemon_ShellServer) error {
	ctx := stream.Context()
	slog := &streamLog{stream: stream, buffered: false}
	log := newStreamLogger(slog)

	stderr := slog.Stderr(false)
	sendErr := func(err error) {
		if errListErr := run.AsErrorList(err); errListErr != nil {
			_ = errListErr.SendToStream(stream)
		} else {
			errStr := err.Error()
			if !strings.HasSuffix(errStr, "\n") {
				errStr += "\n"
			}
			_, _ = slog.Stderr(false).Write([]byte(errStr))
		}
		streamExit(stream, 1)
	}

	ctx, tracer, err := s.beginTracing(ctx, req.AppRoot, req.AppRoot, req.TraceFile)
	if err != nil {
		sendErr(err)
		return nil
	}
	defer func() { _ = tracer.Close() }()

	app, err := s.apps.Track(req.AppRoot)
	if err != nil {
		sendErr(err)
		return nil
	}

	ops := optracker.New(stderr, stream)
	defer ops.AllDone() // Kill the tracker when we exit this function

	p := run.BuildShellParams{
		App:          app,
		BinaryOut:    req.OutputFile,
		WorkingDir:   app.Root(),
		Environ:      req.Environ,
		Stdout:       slog.Stdout(false),
		Stderr:       slog.Stderr(false),
		OpTracker:    ops,
		CodegenDebug: req.CodegenDebug,
		Debug:        req.Debug,
	}
	if buildDir, err := s.mgr.BuildShell(ctx, p); err != nil {
		sendErr(err)
		return nil
	} else {
		if req.CodegenDebug && buildDir != "" {
			log.Info().Msgf("wrote generated code to: %s", buildDir)
		}
		streamExit(stream, 0)
	}

	return nil
}
