//go:build encore_shell

package shellruntime

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"encore.dev/appruntime/shared/appconf"
	"encore.dev/appruntime/shared/logging"
)

var previousLogOut io.Writer = os.Stdout

var rootCmd = &cobra.Command{
	Use:   "",
	Short: fmt.Sprintf("%s Interactive Shell", appconf.Runtime.AppID),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// We override Encore's default logging output to be this shell commands output
		// for this command run, so that we can capture the output and pipe it as it comes
		previousLogOut = logging.LogOut
		logging.LogOut = cmd.OutOrStderr()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		logging.LogOut = previousLogOut
	},
}

func init() {
	rootCmd.AddGroup(
		&cobra.Group{
			ID:    "app_commands",
			Title: fmt.Sprintf("%s Commands", appconf.Runtime.AppID),
		},
		&cobra.Group{
			ID:    "encore_inbuilt",
			Title: "Builtin Commands",
		},
	)
	rootCmd.Version = appconf.Static.AppCommit.Revision
}

// RegisterCommand adds the registered command to the root command
// of the shell.
func RegisterCommand(cmd *cobra.Command) {
	if cmd.GroupID == "" {
		cmd.GroupID = "app_commands"
	}
	rootCmd.AddCommand(cmd)
}
