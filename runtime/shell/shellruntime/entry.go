//go:build encore_shell

package shellruntime

import (
	"fmt"

	shell "github.com/DomBlack/bubble-shell"
	tea "github.com/charmbracelet/bubbletea"

	"encore.dev/appruntime/shared/appconf"
	"encore.dev/shell/shellruntime/internal/tui/helpwrapper"
	"encore.dev/shell/shellruntime/internal/tui/initview"
)

// ShellMain is the entry point for the shell, it is called by the generated main.go
// to start the interactive shell.
func ShellMain() {
	// Startup the interactive shell
	p := tea.NewProgram(tui{view: initview.New()})

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

// tui is the model for the shell's terminal user interface.
//
// It is a simple wrapper over several other models, and is used
// to swap between the init model and the shell model once
// the shell has been initialized.
type tui struct {
	view        tea.Model
	lastSizeMsg tea.WindowSizeMsg
}

var _ tea.Model = tui{}

func (t tui) Init() tea.Cmd {
	return t.view.Init()
}

func (t tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.lastSizeMsg = msg

	case initview.WhoisMsg:
		authedUser = msg

	case initview.AppEnvsMsg:
		knownEnvs = msg.Envs

	case initview.StartShellMsg:
		t.view = helpwrapper.New(
			shell.New(
				rootCmd,
				shell.WithHistoryFile(fmt.Sprintf(".encore-shell/%s.jsonl", appconf.Runtime.AppID)),
			),
		)

		return t, tea.Sequence(
			tea.EnterAltScreen,
			t.view.Init(),
			func() tea.Msg { return t.lastSizeMsg },
		)
	}

	var cmd tea.Cmd
	t.view, cmd = t.view.Update(msg)

	return t, cmd
}

func (t tui) View() string {
	return t.view.View()
}
