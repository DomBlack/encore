package shellruntime

import (
	"io"
	"strings"

	"github.com/spf13/cobra"
)

var (
	traceCmd = &cobra.Command{
		Use:     "trace [on|off]",
		Short:   "Enable or disable API request trace logging",
		GroupID: "encore_inbuilt",
		Aliases: []string{"tracing"},
		Example: `
# View current tracing setting
trace

# Enable trace logging
trace on

# Disable trace logging
trace off
`,
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"on", "off"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				if requestTracingEnabled {
					_, _ = io.WriteString(cmd.OutOrStdout(), "API tracing is currently on\n")
				} else {
					_, _ = io.WriteString(cmd.OutOrStdout(), "API tracing is currently off\n")
				}
			case 1:
				switch strings.ToLower(args[0]) {
				case "on", "true", "yes":
					requestTracingEnabled = true
					_, _ = io.WriteString(cmd.OutOrStdout(), "API tracing is now on\n")
				case "off", "false", "no":
					requestTracingEnabled = false
					_, _ = io.WriteString(cmd.OutOrStdout(), "API tracing is now off\n")
				}
			}
			return nil
		},
	}
	requestTracingEnabled = false
)

func init() {
	rootCmd.AddCommand(traceCmd)
}
