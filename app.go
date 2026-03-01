package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"todo/internal/model"
	"todo/internal/ui"
	"todo/internal/update"
)

// app is a thin wrapper around model.Model that satisfies the tea.Model
// interface by providing Init, Update, and View.
//
// It lives in main rather than in the model package to avoid the import cycle:
//   model → ui → model
//
// Instead the dependency graph is clean:
//   model ← ui
//   model ← update
//   main  → model + ui + update  (only main sees all three)
type app struct {
	m model.Model
}

func (a app) Init() tea.Cmd {
	return a.m.Init()
}

func (a app) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	newModel, cmd := update.Update(a.m, msg)
	return app{m: newModel}, cmd
}

func (a app) View() string {
	return ui.Render(a.m)
}

