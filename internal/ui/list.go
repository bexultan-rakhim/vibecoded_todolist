package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"todo/internal/navigation"
	"todo/internal/task"
)

const (
	// cursorSymbol is the leading glyph on the focused row.
	cursorSymbol    = "󰁔 "
	cursorSymbolNone = "  "
)

// RenderList renders the body panel containing the task list.
// height is the number of rows available for the list body (excluding borders).
// It implements internal scrolling when the task count exceeds the visible area.
func RenderList(tasks []task.Task, cursor navigation.Cursor, width, height int) string {
	if width <= 0 || height <= 0 {
		return ""
	}
	innerWidth := width - 4  // border (2) + padding (2)
	innerHeight := height - 2 // border top + bottom

	var body string
	if len(tasks) == 0 {
		body = renderEmptyState(innerWidth, innerHeight)
	} else {
		body = renderRows(tasks, cursor, innerWidth, innerHeight)
	}

	return ListContainer.
		Width(width - 2).
		Height(height - 2).
		Render(body)
}

// renderRows builds the visible slice of task rows, applying viewport
// scrolling so the focused row is always on screen.
func renderRows(tasks []task.Task, cursor navigation.Cursor, width, height int) string {
	if height <= 0 || width <= 0 {
		return ""
	}
	// Determine the scrolled viewport window.
	start, end := viewportBounds(cursor.Index, len(tasks), height)

	rows := make([]string, 0, end-start)
	for i := start; i < end; i++ {
		rows = append(rows, renderRow(tasks[i], cursor.IsSelected(i), width))
	}

	// Pad with blank lines if the task count is less than the viewport height,
	// so the border stays at a fixed position.
	for len(rows) < height {
		rows = append(rows, strings.Repeat(" ", width))
	}

	return strings.Join(rows, "\n")
}

// viewportBounds returns the [start, end) index range of tasks that should be
// visible given the cursor position and available height.
func viewportBounds(cursorIdx, total, height int) (int, int) {
	if height <= 0 || total <= 0 {
		return 0, 0
	}
	if total <= height {
		return 0, total
	}
	// Keep cursor roughly centred in the viewport.
	half := height / 2
	start := cursorIdx - half
	if start < 0 {
		start = 0
	}
	end := start + height
	if end > total {
		end = total
		start = end - height
	}
	return start, end
}

// renderRow formats a single task row with the appropriate style.
func renderRow(t task.Task, focused bool, width int) string {
// Select the row style first so strikethrough/colour applies uniformly.
	var rowStyle lipgloss.Style
	switch {
	case focused && t.IsDone():
		rowStyle = RowDoneFocused
	case focused:
		rowStyle = RowFocused
	case t.IsDone():
		rowStyle = RowDone
	default:
		rowStyle = RowNormal
	}

	// Plain-text cursor glyph — no styling yet.
	cursorGlyph := cursorSymbolNone
	if focused {
		cursorGlyph = cursorSymbol
	}

	// Plain-text status icon — no styling yet.
	icon := t.Status.Icon()

	// Compute how much space remains for the title after the fixed-width prefix.
	prefixWidth := len([]rune(cursorGlyph)) + len([]rune(icon)) + 1 // +1 for space
	titleWidth := width - prefixWidth - 2                            // -2 for row padding
	title := truncate(t.Title, titleWidth)

	// Compose as plain text, then style everything in one shot.
	content := cursorGlyph + icon + " " + title
	return rowStyle.Width(width).Render(content)
}

// renderEmptyState displays a centred tip when no tasks exist.
func renderEmptyState(width, height int) string {
	msg := EmptyState.Render("No tasks yet. Press 'a' to add one.")
	// Centre vertically.
	topPad := (height - 1) / 2
	lines := make([]string, 0, height)
	for i := 0; i < topPad; i++ {
		lines = append(lines, "")
	}
	lines = append(lines, lipgloss.PlaceHorizontal(width, lipgloss.Center, msg))
	return strings.Join(lines, "\n")
}

// truncate shortens s to at most n visible characters, appending "…" if cut.
func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 1 {
		return "…"
	}
	return string(runes[:n-1]) + "…"
}
