//go:build encore_shell

package shell

import (
	"github.com/spf13/cobra"

	"encore.dev/shell/shellruntime"
)

/*
This file copies in parts of Cobra that are needed for users to write their own
shell command, however it is not used when building the shell itself.
*/

// ShellCompDirective is a bit map representing the different behaviors the shell
// can be instructed to have once completions have been provided.
type ShellCompDirective = cobra.ShellCompDirective

// Command is just that, a command for your application.
// E.g.  'go run ...' - 'run' is the command. Cobra requires
// you to define the usage and description as part of your command
// definition to ensure usability.
type Command = cobra.Command

func register(cmd *Command) {
	shellruntime.RegisterCommand(cmd)
}
