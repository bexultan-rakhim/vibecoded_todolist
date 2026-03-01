package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"todo/config"
	"todo/internal/model"
	"todo/internal/repository"
)

func main() {
	// Resolve the config directory (~/.config/godo or $XDG_CONFIG_HOME/godo).
	dir := config.Dir()

	// Construct the repository — this creates the directory if it doesn't exist.
	repo, err := repository.NewJSONRepository(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "godo: failed to initialise storage: %v\n", err)
		os.Exit(1)
	}

	// Load initial state and construct the model.
	m, err := model.New(repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "godo: failed to load tasks: %v\n", err)
		os.Exit(1)
	}

	// Start the Bubble Tea program.
	// WithAltScreen renders into an alternate terminal buffer so the normal
	// shell output is restored cleanly on exit — no TUI artifacts left behind.
	p := tea.NewProgram(
		app{m: m},
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "godo: %v\n", err)
		os.Exit(1)
	}
}
