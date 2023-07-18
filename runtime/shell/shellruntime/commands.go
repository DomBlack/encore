//go:build encore_shell

package shellruntime

import (
	"fmt"

	"github.com/spf13/cobra"

	"encore.dev/appruntime/shared/appconf"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: fmt.Sprintf("%s Interactive Shell", appconf.Runtime.AppID),
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
