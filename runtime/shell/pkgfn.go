package shell

// RegisterCommand registers a command with the interactive shell.
//
// This commands registered here will have no effect on your normal
// application. However when you run `encore shell` they will be
// compiled into the interactive shell and will be available to the
// user to execute.
func RegisterCommand(cmd *Command) {
	register(cmd)
}
