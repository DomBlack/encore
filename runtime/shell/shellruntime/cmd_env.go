//go:build encore_shell

package shellruntime

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var (
	envCmd = &cobra.Command{
		Use:     "env [env name]",
		Short:   "Change the current environment being access",
		GroupID: "encore_inbuilt",
		Aliases: []string{"environment"},
		Example: `
# List all available environments
env

# Change to your local development environment
env local

# Change to your cloud based "staging" environment
env staging
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return listEnvs(cmd.OutOrStdout())
			}
			return changeEnv(cmd.OutOrStdout(), args[0])
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			switch len(args) {
			case 0: // the first argument is the environment name
				return envAutoComplete(toComplete), cobra.ShellCompDirectiveNoFileComp
			default:
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
		},
	}
	activeEnv = "local"                 // activeEnv is the currently active environment
	knownEnvs = make(map[string]string) // knownEnvs is a map of all known environments
)

func init() {
	rootCmd.AddCommand(envCmd)

	if err := changeEnv(io.Discard, "local"); err != nil {
		panic(fmt.Sprintf("failed to set initial environment: %s", err))
	}
}

func listEnvs(out io.Writer) error {
	// Create a sorted list of all known environments
	envs := make([]string, 0, len(knownEnvs)+1)
	envs = append(envs, "local")
	for env := range knownEnvs {
		envs = append(envs, env)
	}
	sort.Strings(envs)

	// Print out all known environments
	_, _ = io.WriteString(out, "Environments:\n")
	for _, env := range envs {
		_, _ = io.WriteString(out, "  ")
		if env == activeEnv {
			_, _ = io.WriteString(out, "* ")
		} else {
			_, _ = io.WriteString(out, "  ")
		}
		_, _ = io.WriteString(out, env)
		_, _ = io.WriteString(out, "\n")
	}

	return nil
}

func changeEnv(out io.Writer, newEnv string) error {
	if newEnv != "local" && knownEnvs[newEnv] == "" {
		return fmt.Errorf("unknown environment: %s", newEnv)
	}
	activeEnv = newEnv

	// Update the active transport
	switch activeEnv {
	case "local":
		shellApiTransportSingleton.ActiveTransport = &localLoopBackTransport{}
	default:
		shellApiTransportSingleton.ActiveTransport = &encorePlatformProxyTransport{
			envName: newEnv,
		}
	}

	_, _ = io.WriteString(out, "Environment changed to ")
	_, _ = io.WriteString(out, newEnv)
	_, _ = io.WriteString(out, "\n")

	return nil
}

func envAutoComplete(toComplete string) []string {
	envs := make([]string, 0, len(knownEnvs)+1)

	if toComplete == "" || strings.HasPrefix("local", toComplete) {
		envs = append(envs, "local")
	}
	for env := range knownEnvs {
		if toComplete == "" || strings.HasPrefix(env, toComplete) {
			envs = append(envs, env)
		}
	}
	sort.Strings(envs)

	return envs
}
