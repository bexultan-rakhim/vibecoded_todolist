package ui

import (
	"github.com/charmbracelet/lipgloss"
	"todo/internal/model"
)

const (
	// headerHeight is the fixed height of the header panel (border + 1 content line).
	headerHeight = 3
	// footerHeight is the fixed height of the footer panel (border + 1 content line).
	footerHeight = 3
)

// Render is the top-level View function for Godo.
// It composes the three layout zones — Header, Body, Footer — into a single
// string that Bubble Tea writes to the terminal on every update cycle.
//
// The body height is derived dynamically from the terminal height so that the
// header and footer remain pinned while the list scrolls internally.
func Render(m model.Model) string {
	w := m.Width
	h := m.Height

	// Guard against zero dimensions during startup before the first
	// tea.WindowSizeMsg has been received.
	if w == 0 || h == 0 {
		return ""
	}

	bodyHeight := h - headerHeight - footerHeight

	header := RenderHeader(m.Tasks, w)
	body := RenderList(m.Tasks, m.Cursor, w, bodyHeight)
	footer := RenderFooter(m, w)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}
