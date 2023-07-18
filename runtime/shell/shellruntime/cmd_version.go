//go:build encore_shell

package shellruntime

import (
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"

	"encore.dev/appruntime/shared/appconf"
)

var (
	versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Print the version information for this shell",
		GroupID: "encore_inbuilt",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return printVersion(cmd.OutOrStdout())
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(out io.Writer) error {
	_, _ = io.WriteString(out, fmt.Sprintf("Shell:   %s/%s\n", appconf.Runtime.AppID, appconf.Static.AppCommit.Revision))
	_, _ = io.WriteString(out, fmt.Sprintf("Runtime: %s\n", appconf.Static.EncoreCompiler))
	_, _ = io.WriteString(out, fmt.Sprintf("Go:      %s\n", runtime.Version()))
	_, _ = io.WriteString(out, "\n")

	return nil
}
