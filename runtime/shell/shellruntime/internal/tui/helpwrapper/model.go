//go:build encore_shell

package helpwrapper

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// helpWrapperModel is a single model in which we've embedded the shell in another model
// which also displays help
type helpWrapperModel struct {
	Shell tea.Model
	Help  help.Model

	width     int
	helpStyle lipgloss.Style
}

// New creates a new help wrapper model which renders the given model with a help footer
func New(shell tea.Model) tea.Model {
	helpModel := help.New()
	helpModel.Styles.ShortKey = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))

	return helpWrapperModel{
		Shell: shell,
		Help:  helpModel,
		helpStyle: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder(), true, false, false).
			BorderForeground(lipgloss.Color("#4E4E4E")).
			Align(lipgloss.Right),
	}
}

var _ tea.Model = helpWrapperModel{}

func (m helpWrapperModel) Init() tea.Cmd {
	return m.Shell.Init()
}

func (m helpWrapperModel) Update(msg tea.Msg) (rtn tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width

		m.Shell, cmd = m.Shell.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - 2, // Leave space for the help and it's border
		})
	default:
		m.Shell, cmd = m.Shell.Update(msg)
	}

	return m, cmd
}

func (m helpWrapperModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		m.Shell.View(),
		m.helpStyle.Copy().Width(m.width).Render(
			m.Help.View(
				m.Shell.(help.KeyMap), // The shell implements the help.KeyMap interface
			),
		),
	)
}
