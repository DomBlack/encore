//go:build encore_shell

package shellruntime

import (
	"fmt"

	shell "github.com/DomBlack/bubble-shell"
	tea "github.com/charmbracelet/bubbletea"

	"encore.dev/appruntime/shared/appconf"
	"encore.dev/shell/shellruntime/internal/tui/helpwrapper"
)

// ShellMain is the entry point for the shell, it is called by the generated main.go
// to start the interactive shell.
func ShellMain() {
	// Startup the interactive shell
	appName := appconf.Runtime.AppID
	p := tea.NewProgram(
		helpwrapper.New(
			shell.New(
				rootCmd,
				shell.WithHistoryFile(fmt.Sprintf(".encore-shell/%s.jsonl", appName)),
			),
		),
	)

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
