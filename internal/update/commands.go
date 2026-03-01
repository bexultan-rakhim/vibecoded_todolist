// Package update contains the Bubble Tea Update function and Cmd factories.
// Commands are the only mechanism for side effects in the MVU architecture —
// they run outside the update cycle and deliver results back as messages.
package update

import (
	tea "github.com/charmbracelet/bubbletea"
	"todo/internal/repository"
	"todo/internal/task"
)

// --- Messages ---

// SavedMsg is delivered when a Save command completes successfully.
type SavedMsg struct{}

// SaveErrMsg is delivered when a Save command fails.
type SaveErrMsg struct {
	Err error
}

// --- Commands ---

// SaveCmd returns a Cmd that persists the given task list via the repository.
// On success it emits SavedMsg; on failure it emits SaveErrMsg.
// It is called after every mutation so the acceptance criterion
// ("all changes saved immediately upon mutation") is met.
func SaveCmd(repo repository.Repository, tasks []task.Task) tea.Cmd {
	return func() tea.Msg {
		if err := repo.Save(tasks); err != nil {
			return SaveErrMsg{Err: err}
		}
		return SavedMsg{}
	}
}
