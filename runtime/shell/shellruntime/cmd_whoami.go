//go:build encore_shell

package shellruntime

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"encore.dev/shell/shellruntime/internal/tui/initview"
)

var (
	whoamiCmd = &cobra.Command{
		Use:     "whoami",
		Short:   "Print information about the current user",
		GroupID: "encore_inbuilt",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return whoami(cmd.OutOrStdout())
		},
	}
	authedUser initview.WhoisMsg
)

func init() {
	rootCmd.AddCommand(whoamiCmd)
}

func whoami(out io.Writer) error {
	if authedUser.ID == "" {
		return fmt.Errorf("not logged in")
	}

	_, _ = fmt.Fprintf(out,
		"Logged in as `%s`\n",
		lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(authedUser.Email),
	)

	return nil
}
