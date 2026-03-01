package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"todo/internal/task"
)

// RenderHeader renders the top panel containing the app title and task statistics.
// It adapts its width to the available terminal width passed in from the model.
func RenderHeader(tasks []task.Task, width int) string {
	total := len(tasks)
	done := 0
	for _, t := range tasks {
		if t.IsDone() {
			done++
		}
	}
	todo := total - done

	title := HeaderTitle.Render("󰄱 godo")

	stats := lipgloss.JoinHorizontal(
		lipgloss.Left,
		HeaderStatsDone.Render(fmt.Sprintf("󰄵 %d done", done)),
		HeaderStats.Render("  ·  "),
		HeaderStats.Render(fmt.Sprintf("󰄱 %d todo", todo)),
	)

	// Push stats to the right side of the header.
	innerWidth := max(0, width-4) // account for border (2) + padding (2)
	gap := max(0, innerWidth-lipgloss.Width(title)-lipgloss.Width(stats))
	spacer := lipgloss.NewStyle().Render(repeatSpace(gap))

	inner := lipgloss.JoinHorizontal(lipgloss.Left, title, spacer, stats)

	return HeaderContainer.
		Width(width - 2). // subtract border width
		Render(inner)
}

func repeatSpace(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := range b {
		b[i] = ' '
	}
	return string(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
