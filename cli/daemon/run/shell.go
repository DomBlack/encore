package run

import (
	"context"
	"fmt"
	"io"
	"net/netip"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/rs/zerolog/log"

	encore "encore.dev"
	"encore.dev/appruntime/exported/experiments"
	"encr.dev/cli/daemon/apps"
	"encr.dev/internal/optracker"
	"encr.dev/pkg/builder"
	"encr.dev/pkg/builder/builderimpl"
	"encr.dev/pkg/cueutil"
	"encr.dev/pkg/option"
	"encr.dev/pkg/paths"
	"encr.dev/pkg/svcproxy"
	"encr.dev/pkg/vcs"
	metav1 "encr.dev/proto/encore/parser/meta/v1"
)

// BuildShellParams groups the parameters for the Shell method.
type BuildShellParams struct {
	// App is the app to execute the script for.
	App *apps.Instance

	// BinaryOut is the path to write the binary to.
	BinaryOut string

	// WorkingDir is the working dir to execute the script from.
	// It's relative to the app root.
	WorkingDir string

	// Environ are the environment variables to set when running the tests,
	// in the same format as os.Environ().
	Environ []string

	// Stdout and Stderr are where "go test" output should be written.
	Stdout, Stderr io.Writer

	OpTracker *optracker.OpTracker

	// Debug specifies to compile the application for debugging.
	Debug bool

	// CodegenDebug, if true, specifies to keep the output
	// around for codegen debugging purposes.
	CodegenDebug bool
}

func (mgr *Manager) BuildShell(ctx context.Context, p BuildShellParams) (buildDir string, err error) {
	// Return early if the ctx is already canceled.
	if err := ctx.Err(); err != nil {
		return "", err
	}

	tracker := p.OpTracker
	jobs := optracker.NewAsyncBuildJobs(ctx, p.App.PlatformOrLocalID(), tracker)

	// Parse the app to figure out what infrastructure is needed.
	start := time.Now()
	parseOp := tracker.Add("Building Encore application graph", start)
	topoOp := tracker.Add("Analyzing service topology", start)

	expSet, err := p.App.Experiments(p.Environ)
	if err != nil {
		return "", errors.Wrap(err, "get experimental features")
	}
	expSet.Add(experiments.ExternalCalls)

	vcsRevision := vcs.GetRevision(p.App.Root())
	buildInfo := builder.BuildInfo{
		BuildTags:          builder.ShellBuildTags,
		CgoEnabled:         true,
		StaticLink:         false,
		Debug:              p.Debug,
		GOOS:               runtime.GOOS,
		GOARCH:             runtime.GOARCH,
		KeepOutput:         p.CodegenDebug,
		Revision:           vcsRevision.Revision,
		UncommittedChanges: vcsRevision.Uncommitted,
		MainPkg:            option.Some[paths.Pkg]("./__encore/shell"),
		BuildShell:         true,
	}

	bld := builderimpl.Resolve(expSet)
	parse, err := bld.Parse(ctx, builder.ParseParams{
		Build:       buildInfo,
		App:         p.App,
		Experiments: expSet,
		WorkingDir:  p.App.Root(),
		ParseTests:  false,
	})
	if err != nil {
		tracker.Fail(parseOp, err)
		return "", err
	}
	tracker.Done(parseOp, 500*time.Millisecond)
	tracker.Done(topoOp, 300*time.Millisecond)

	// Build a basic runtime config for the shell to embed.
	apiBaseURL := fmt.Sprintf("http://localhost:%d", mgr.RuntimePort)
	svcProxy, err := svcproxy.New(ctx, log.Logger)
	if err != nil {
		return "", errors.Wrap(err, "create service proxy")
	}
	defer svcProxy.Close()
	listenAddresses, err := GenerateListenAddresses(svcProxy, parse.Meta.Svcs)
	if err != nil {
		return "", errors.Wrap(err, "generate listen addresses")
	}

	envGen := &RuntimeEnvGenerator{
		App:             p.App,
		InfraManager:    nil,
		Meta:            parse.Meta,
		Secrets:         nil,
		SvcConfigs:      nil,
		DaemonProxyAddr: option.Some(netip.AddrPortFrom(netip.IPv4Unspecified(), uint16(mgr.RuntimePort))),
		ListenAddresses: listenAddresses,
		AppID:           option.Some(p.App.PlatformID()),
		EnvID:           option.Some("shell"),
		EnvName:         option.Some("shell"),
		EnvType:         option.Some(encore.EnvDevelopment),
		CloudType:       option.Some(encore.CloudLocal),
		ServiceAuthType: option.Some("noop"),
	}
	runtimeCfg, err := envGen.runtimeConfigForServices([]*metav1.Service{
		{ // we add a fake service here so all the "IsHosting" code and "Gateway" checks
			// return false
			Name:    "__encore_shell",
			RelPath: "__encore/shell",
		},
	}, false)
	if err != nil {
		return "", errors.Wrap(err, "generate runtime config")
	}
	buildInfo.ShellEnvs = []string{
		runtimeCfgEnvVar + "=" + runtimeCfg,
	}

	// Now compile the application.
	var build *builder.CompileResult
	jobs.Go("Compiling application source code", false, 0, func(ctx context.Context) (err error) {
		build, err = bld.Compile(ctx, builder.CompileParams{
			Build:       buildInfo,
			App:         p.App,
			Parse:       parse,
			OpTracker:   tracker,
			Experiments: expSet,
			WorkingDir:  p.WorkingDir,
			CueMeta: &cueutil.Meta{
				APIBaseURL: apiBaseURL,
				EnvName:    "local",
				EnvType:    cueutil.EnvType_Development,
				CloudType:  cueutil.CloudType_Local,
			},
		})
		if err != nil {
			return errors.Wrap(err, "compile error on exec")
		}
		return nil
	})
	defer func() {
		if build != nil && build.Dir != "" {
			if p.CodegenDebug {
				buildDir = build.Dir
			} else {
				_ = os.RemoveAll(build.Dir)
			}
		}
	}()

	if err := jobs.Wait(); err != nil {
		return "", err
	}

	// Copy the build.exe to requested location
	if err := copyFile(build.Exe, p.BinaryOut); err != nil {
		return "", errors.Wrap(err, "copy shell binary")
	}

	tracker.AllDone()
	return "", nil
}

func copyFile(from, to string) error {
	// Ensure the destination directory exists
	if err := os.MkdirAll(filepath.Dir(to), 0755); err != nil {
		return errors.Wrap(err, "create destination directory")
	}

	// Create the src/dst file objects
	src, err := os.Open(from)
	if err != nil {
		return errors.Wrapf(err, "open source file (%s)", from)
	}
	defer func() { _ = src.Close() }()

	dst, err := os.Create(to)
	if err != nil {
		return errors.Wrapf(err, "create destination file (%s)", to)
	}
	defer func() { _ = dst.Close() }()

	// Read the permissions from the src file
	fi, err := src.Stat()
	if err != nil {
		return errors.Wrap(err, "stat source file")
	}

	// Copy the contents of the src file to the dst file
	if _, err := io.Copy(dst, src); err != nil {
		return errors.Wrap(err, "copy file")
	}

	if err := dst.Close(); err != nil {
		return errors.Wrap(err, "close destination file")
	}

	// Set the permissions on the dst file
	if err := os.Chmod(to, fi.Mode()); err != nil {
		return errors.Wrap(err, "chmod destination file")
	}

	return nil
}
