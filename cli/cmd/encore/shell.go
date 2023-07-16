package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"encr.dev/cli/cmd/encore/cmdutil"
	"encr.dev/cli/cmd/encore/root"
	daemonpb "encr.dev/proto/encore/daemon"
)

var shellCmd = &cobra.Command{
	Use:   "shell [--env=local] [--debug] [--codegen-debug]",
	Short: "Start an interactive shell",
	Example: `
# Start an interactive shell against the locally running app 
encore shell

# Start an interactive shell against your "prod" environment
encore shell --env=prod
`,
	Run: func(cmd *cobra.Command, args []string) {
		appRoot, _ := determineAppRoot()

		runAfterBuild := shellOutputPath == ""
		if runAfterBuild {
			// Create a temporary directory to build the shell into
			tmp, err := os.MkdirTemp("", "encore-shell-")
			if err != nil {
				cmdutil.Fatalf("unable to create temporary directory: %v", err)
			}
			shellOutputPath = filepath.Join(tmp, "shell")

			// After we've exited, remove the temporary directory
			defer func() { _ = os.RemoveAll(filepath.Dir(shellOutputPath)) }()
		}

		// Build the shell
		cmd.Context()
		daemon := setupDaemon(cmd.Context())
		stream, err := daemon.Shell(cmd.Context(), &daemonpb.ShellRequest{
			AppRoot:      appRoot,
			OutputFile:   shellOutputPath,
			Debug:        debug,
			CodegenDebug: codegenDebug,
			Environ:      os.Environ(),
			TraceFile:    root.TraceFile,
		})
		if err != nil {
			cmdutil.Fatalf("unable to build shell: %v", err)
		}

		// Stream the build output to the terminal
		exitCodeFromBuild := streamCommandOutput(stream, convertJSONLogs())
		if exitCodeFromBuild != 0 {
			// If the build failed, exit with the same code
			os.Exit(exitCodeFromBuild)
		}

		// If the user just wanted the binary, we're done
		if !runAfterBuild {
			os.Exit(0)
		}

		// Otherwise start the shell up with the args passed in
		shellCmd := exec.CommandContext(cmd.Context(), shellOutputPath, args...)
		shellCmd.Stdin = os.Stdin
		shellCmd.Stdout = os.Stdout
		shellCmd.Stderr = os.Stderr
		if err := shellCmd.Run(); err != nil {
			cmdutil.Fatalf("unable to run shell: %v", err)
		}
		os.Exit(0)
	},
}

var (
	shellEnv        string
	shellOutputPath string
)

func init() {
	shellCmd.Flags().StringVarP(&shellEnv, "env", "e", "local", "Environment to connect to")
	shellCmd.Flags().StringVarP(&shellOutputPath, "output", "o", "", "If set, the shell will be built to this path instead of run")
	shellCmd.Flags().BoolVar(&codegenDebug, "codegen-debug", false, "Dump generated code (for debugging Encore's code generation)")
	shellCmd.Flags().BoolVar(&debug, "debug", false, "Compile for debugging (disables some optimizations)")
	alphaCmd.AddCommand(shellCmd)
}
