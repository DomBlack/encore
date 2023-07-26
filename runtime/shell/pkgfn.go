package shell

// NewCommand creates a command within the interactive shell of this application.
//
// This commands registered here will have no effect on your normal
// application. However when you run `encore shell` they will be
// compiled into the interactive shell and will be available to the
// user to execute.
func NewCommand(cmd *Command) *RegisteredCommand {
	register(cmd)
	return &RegisteredCommand{}
}

// RegisteredCommand is a command that has been registered with the applications
// shell using the [NewCommand] function.
type RegisteredCommand struct{}
