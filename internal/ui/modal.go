package ui

import (
	"github.com/charmbracelet/lipgloss"
	"todo/internal/model"
)

const (
	// modalWidth is the fixed inner width of the modal box.
	// Wide enough for comfortable typing, narrow enough to feel focused.
	modalWidth = 50
)

// RenderModal returns a centered overlay string for ModeAdd and ModeEdit.
// It composes the modal box over a dimmed version of the existing background
// content, then centres it within the terminal dimensions.
//
// bgContent is the already-rendered base layout (header + list + footer)
// which is displayed dimmed behind the modal.
func RenderModal(m model.Model, bgContent string) string {
	if m.Mode != model.ModeAdd && m.Mode != model.ModeEdit {
		return bgContent
	}

	modal := renderModalBox(m)

	// Place the modal in the centre of the terminal.
	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(colorFgMuted),
	)
}

// renderModalBox builds the bordered input box content.
func renderModalBox(m model.Model) string {
	title := modalTitle(m.Mode)
	hint := ModalHint.Render("enter to confirm · esc to cancel")
	input := m.Input.View()

	// Stack: title → input → hint
	inner := lipgloss.JoinVertical(
		lipgloss.Left,
		ModalTitle.Render(title),
		input,
		"",
		hint,
	)

	return ModalContainer.
		Width(modalWidth).
		Render(inner)
}

// modalTitle returns the heading text for the modal depending on the active mode.
func modalTitle(mode model.ViewMode) string {
	switch mode {
	case model.ModeAdd:
		return "󰐕  New Task"
	case model.ModeEdit:
		return "  Edit Task"
	default:
		return ""
	}
}
