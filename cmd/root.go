package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Camilo-845/type-cli/internal/tui"
)

var rootCmd = &cobra.Command{
	Use:   "tcli",
	Short: "A Monkeytype-inspired typing game for the terminal",
	Long: `TypeGame is a terminal-based typing speed test inspired by Monkeytype.
Test your typing speed with different modes, word lists, and track your progress.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
