package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/Camilo-845/type-cli/internal/tui"
	"github.com/Camilo-845/type-cli/internal/words"
)

var languageFlag string

var rootCmd = &cobra.Command{
	Use:   "tcli",
	Short: "A Monkeytype-inspired typing game for the terminal",
	Long: `TypeGame is a terminal-based typing speed test inspired by Monkeytype.
Test your typing speed with different modes, word lists, and track your progress.`,
	Run: func(cmd *cobra.Command, args []string) {
		model := tui.NewModel()

		if languageFlag != "" {
			available := make(map[string]bool)
			for _, lang := range words.ListLanguages() {
				available[lang] = true
			}
			if !available[languageFlag] {
				fmt.Fprintf(os.Stderr, "unknown language: %s\n", languageFlag)
				fmt.Fprintf(os.Stderr, "available: %v\n", words.ListLanguages())
				os.Exit(1)
			}
			model.SetLanguage(languageFlag)
		}

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringVarP(&languageFlag, "language", "l", "", "language to use (e.g. english, spanish, english_1k)")
	rootCmd.RegisterFlagCompletionFunc("language", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return words.ListLanguages(), cobra.ShellCompDirectiveNoFileComp
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
