//go:build encore_shell

package initview

import (
	"context"
	"fmt"
	"time"

	"github.com/DomBlack/bubble-shell/pkg/tui/errdisplay"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cockroachdb/errors"

	"encore.dev/appruntime/shared/appconf"
	"encore.dev/shell/shellruntime/internal/platform"
	"encore.dev/shell/shellruntime/internal/platform/auth"
)

type Model struct {
	ctx         context.Context
	ctxCancel   context.CancelFunc
	spinner     spinner.Model
	errDisplay  tea.Model
	loginNeeded bool
	tasks       []*Task
}

var _ help.KeyMap = Model{}

func New() tea.Model {
	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		ctx:       ctx,
		ctxCancel: cancel,
		spinner: spinner.New(spinner.WithSpinner(spinner.Spinner{
			Frames: []string{"⠋", "⠙", "⠚", "⠒", "⠂", "⠂", "⠒", "⠲", "⠴", "⠦", "⠖", "⠒", "⠐", "⠐", "⠒", "⠓", "⠋"},
			FPS:    time.Second / 4,
		})),
		errDisplay: nil,
		tasks: []*Task{
			{
				Description: "Authenticating with Encore",
				Status:      Running,
			},
			{
				Description: "Getting environment list",
			},
			{
				Description: "Verifying permissions",
			},
			{
				Description: "Starting shell",
			},
		},
	}
}

var cancel = key.NewBinding(
	key.WithKeys("ctrl+c"),
	key.WithHelp("ctrl+c", "Cancel"),
)

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.InitPlatformClient,
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case WhoisMsg:
		m.tasks[0].Status = Success
		m.tasks[0].FinishMsg = fmt.Sprintf(
			"Authenticated as `%s%s",
			lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Render(msg.Email),
			lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("`"), // This renders the closing backtick as green (matching the Success color)
		)
		m.tasks[1].Status = Running

		return m, m.GetAppEnvironments

	case AppEnvsMsg:
		m.tasks[1].Status = Success
		m.tasks[2].Status = Running

		return m, m.VerifyPermissions

	case PermissionsVerifiedMsg:
		m.tasks[2].Status = Success
		m.tasks[3].Status = Running

		return m, m.StartShell

	case StartShellMsg:
		m.tasks[3].Status = Success

		// No-op here, we've initialized everything we need to
		// and can drop into the interactive shell
		return m, nil

	case errorMsg:
		for _, task := range m.tasks {
			if task.Status == Running {
				task.Status = Failure
			} else if task.Status != Success && task.Status != Failure {
				task.Status = Skipped
			}
		}

		if errors.Is(msg.err, auth.ErrNotLoggedIn) {
			m.loginNeeded = true

			// We wait a second which gives the user a chance to see the error
			// before we quit
			return m, func() tea.Msg {
				time.Sleep(time.Second)
				return tea.Quit()
			}
		} else {
			// If we're not already displaying an error, display this one
			if m.errDisplay == nil {
				m.errDisplay = errdisplay.New(nil, msg.err)
				return m, m.errDisplay.Init()
			}
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, cancel):
			m.ctxCancel()
			return m, tea.Quit
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	if m.errDisplay != nil {
		m.errDisplay, cmd = m.errDisplay.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var lines []string

	// Render the tasks
	for _, item := range m.tasks {
		status, suffix, itemColor := m.itemRender(item)

		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left,
			"  ",
			lipgloss.NewStyle().Foreground(itemColor).Width(2).Render(status),
			lipgloss.NewStyle().Foreground(itemColor).Render(item.Description+suffix),
		))
	}

	// If we have an error display, add it to the end
	if m.loginNeeded {
		red := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
		cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

		lines = append(lines,
			"",
			red.Render("You need to be logged into Encore to use the shell. Run `")+cyan.Render("encore auth login")+red.Render("` on this machine to login, then try again."),
			"",
			"",
		)
	} else if m.errDisplay != nil {
		lines = append(lines,
			"",
			m.errDisplay.View(),
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("Press ctrl+c to quit"),
		)
	} else {
		lines = append(lines,
			"",
			lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Render("Press ctrl+c to cancel"),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Top, lines...)
}

// InitPlatformClient creates a platform client and checks we're logged in
func (m Model) InitPlatformClient() (rtn tea.Msg) {
	defer func() {
		if err := recover(); err != nil {
			rtn = errorMsg{errors.Newf("panic: %v", err)}
		}
	}()

	// Create a client
	if err := platform.Init(); err != nil {
		return errorMsg{err}
	}

	// Check we're logged in
	qry := `query WhoAmI {
  me {
    id
    email
    full_name
  }
}`
	type WhoAmIResp struct {
		Me struct {
			ID       string `json:"id"`
			Email    string `json:"email"`
			FullName string `json:"full_name"`
		} `json:"me"`
	}
	resp := &WhoAmIResp{}
	err := platform.Query(m.ctx, qry, nil, resp)
	if err != nil {
		return errorMsg{err}
	}

	return WhoisMsg{
		ID:    resp.Me.ID,
		Email: resp.Me.Email,
		Name:  resp.Me.FullName,
	}
}

func (m Model) GetAppEnvironments() (rtn tea.Msg) {
	defer func() {
		if err := recover(); err != nil {
			rtn = errors.Newf("panic: %v", err)
		}
	}()

	// Query for the app environments
	qry := `query GetEnvs($appSlug:String!) {
  app(slug:$appSlug){
    id
    envs {
      id
      name      
    }
  }
}`
	type GetEnvsResp struct {
		App struct {
			ID   string `json:"id"`
			Envs []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			}
		}
	}
	resp := &GetEnvsResp{}
	err := platform.Query(m.ctx, qry, map[string]any{
		"appSlug": appconf.Runtime.AppID,
	}, resp)
	if err != nil {
		return errorMsg{err}
	}

	// Build the message from the response
	msg := AppEnvsMsg{Envs: make(map[string]string)}
	for _, env := range resp.App.Envs {
		msg.Envs[env.Name] = env.ID
	}
	return msg
}

func (m Model) VerifyPermissions() tea.Msg {
	time.Sleep(222 * time.Millisecond)

	return PermissionsVerifiedMsg{}
}

func (m Model) StartShell() tea.Msg {
	time.Sleep(333 * time.Millisecond)

	return StartShellMsg{}
}

func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{cancel}
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{{cancel}}
}

type Status uint8

const (
	Waiting Status = iota
	Running
	Success
	Skipped
	Failure
)

type Task struct {
	Description string
	FinishMsg   string
	Status      Status
}

func (m Model) itemRender(t *Task) (status, suffix string, itemColor lipgloss.Color) {
	itemColor = "#4A4A4A"

	// These colours are based on Encore's own internal tracker
	// this means the tasks we render when the shell starts up
	// matches the colors of the original compiler tasks
	switch t.Status {
	case Waiting:
		status = "⏳"
		itemColor = "6"
	case Running:
		status = m.spinner.View()
		suffix = "..."
		// Cyan
		itemColor = "6"
	case Failure:
		status = "❌"
		suffix = "... Failed"
		// Red
		itemColor = "1"
	case Success:
		status = "✔"
		if t.FinishMsg != "" {
			suffix = "... " + t.FinishMsg
		} else {
			suffix = "... Done!"
		}
		// Green
		itemColor = "2"
	case Skipped:
		status = "⚠"
		suffix = "... Canceled"
		// Yellow
		itemColor = "3"
	}

	return
}

type errorMsg struct {
	err error
}

// WhoisMsg is a message that contains the user's ID, email and name
// once they've been authenticated against the Encore platform.
type WhoisMsg struct {
	ID    string
	Email string
	Name  string
}

// AppEnvsMsg contains a map of environment names to IDs
// as returned by the Encore platform.
type AppEnvsMsg struct {
	Envs map[string]string // Name of environment name to ID
}

// PermissionsVerifiedMsg is a message that indicates the user's
// permissions have been verified against the given Environment.
type PermissionsVerifiedMsg struct {
}

// StartShellMsg is a message that we're ready to drop
// the user into the shell.
type StartShellMsg struct {
}
