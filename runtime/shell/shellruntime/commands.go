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

// RegisterCommand adds the registered command to the root command
// of the shell.
func RegisterCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
