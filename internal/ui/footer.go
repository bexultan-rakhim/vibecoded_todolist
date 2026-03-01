package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"todo/internal/model"
)

// hint is a key+description pair for the footer.
type hint struct {
	key  string
	desc string
}

// listModeHints are the hints shown in ModeList.
var listModeHints = []hint{
	{"j/k", "navigate"},
	{"a", "add"},
	{"e", "edit"},
	{"x", "delete"},
	{"enter", "toggle"},
	{"u", "undo"},
	{"r", "redo"},
	{"q", "quit"},
}

// inputModeHints are the hints shown in ModeAdd and ModeEdit.
var inputModeHints = []hint{
	{"enter", "confirm"},
	{"esc", "cancel"},
}

// RenderFooter renders the bottom panel.
// If the model carries a non-nil error, it displays the error instead of hints.
func RenderFooter(m model.Model, width int) string {
	var inner string

	if m.Err != nil {
		inner = FooterError.Render("error: " + m.Err.Error())
	} else {
		hints := listModeHints
		if m.Mode == model.ModeAdd || m.Mode == model.ModeEdit {
			hints = inputModeHints
		}
		inner = renderHints(hints, m)
	}

	return FooterContainer.
		Width(width - 2).
		Render(inner)
}

// renderHints builds the space-separated key hint string, dimming undo/redo
// hints when the history stack has nothing to act on.
func renderHints(hints []hint, m model.Model) string {
	sep := FooterSeparator.Render(" · ")
	parts := make([]string, 0, len(hints))

	for _, h := range hints {
		key := KeyHintKey.Render(h.key)
		desc := KeyHintDesc.Render(h.desc)

		// Dim undo/redo when they are not actionable.
		if h.key == "u" && !m.History.CanUndo() {
			key = KeyHintKey.Copy().Foreground(colorFgMuted).Render(h.key)
			desc = KeyHintDesc.Copy().Foreground(colorBorder).Render(h.desc)
		}
		if h.key == "r" && !m.History.CanRedo() {
			key = KeyHintKey.Copy().Foreground(colorFgMuted).Render(h.key)
			desc = KeyHintDesc.Copy().Foreground(colorBorder).Render(h.desc)
		}

		parts = append(parts, lipgloss.JoinHorizontal(lipgloss.Left, key, " ", desc))
	}

	return strings.Join(parts, sep)
}
