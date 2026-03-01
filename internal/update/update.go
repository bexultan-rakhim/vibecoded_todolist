package update

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"todo/internal/history"
	"todo/internal/model"
	"todo/internal/task"
)

// Update is the Bubble Tea Update function for Godo.
// It is a pure function: given the current model and an incoming message,
// it returns the next model and an optional side-effect command.
// It is split into sub-handlers per mode to keep each section focused.
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Terminal resize — update dimensions so View re-renders responsively.
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	// Save completed successfully — clear any displayed error.
	case SavedMsg:
		m.Err = nil
		return m, nil

	// Save failed — surface the error in the footer.
	case SaveErrMsg:
		m.Err = msg.Err
		return m, nil

	case tea.KeyMsg:
		switch m.Mode {
		case model.ModeList:
			return updateList(m, msg)
		case model.ModeAdd, model.ModeEdit:
			return updateInput(m, msg)
		}
	}

	return m, nil
}

// updateList handles all keypresses in ModeList.
func updateList(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	km := model.DefaultKeyMap()

	switch {

	// --- Navigation ---
	case key.Matches(msg, km.Up):
		m.Cursor = m.Cursor.Up(len(m.Tasks))

	case key.Matches(msg, km.Down):
		m.Cursor = m.Cursor.Down(len(m.Tasks))

	case key.Matches(msg, km.Top):
		m.Cursor = m.Cursor.Top(len(m.Tasks))

	case key.Matches(msg, km.Bottom):
		m.Cursor = m.Cursor.Bottom(len(m.Tasks))

	// --- Add ---
	case key.Matches(msg, km.Add):
		m.Input.SetValue("")
		m.Input.Focus()
		m.Mode = model.ModeAdd

	// --- Edit ---
	case key.Matches(msg, km.Edit):
		if len(m.Tasks) == 0 {
			break
		}
		m.Input.SetValue(m.Tasks[m.Cursor.Index].Title)
		m.Input.Focus()
		m.Mode = model.ModeEdit

	// --- Toggle done/todo ---
	case key.Matches(msg, km.Toggle):
		if len(m.Tasks) == 0 {
			break
		}
		m.Tasks[m.Cursor.Index].Toggle()
		return m, SaveCmd(m.Repo, m.Tasks)

	// --- Delete ---
	case key.Matches(msg, km.Delete):
		if len(m.Tasks) == 0 {
			break
		}
		idx := m.Cursor.Index
		// Record snapshot before deletion for undo.
		m.History.Push(history.Entry{
			Kind:     history.ActionDelete,
			Snapshot: m.Tasks[idx],
			Index:    idx,
		})
		m.Tasks = removeTask(m.Tasks, idx)
		m.Cursor = m.Cursor.Clamp(len(m.Tasks))
		return m, SaveCmd(m.Repo, m.Tasks)

	// --- Undo ---
	case msg.String() == "u":
		m = applyUndo(m)
		if m.Err == nil {
			return m, SaveCmd(m.Repo, m.Tasks)
		}

	// --- Redo ---
	case msg.String() == "r":
		m = applyRedo(m)
		if m.Err == nil {
			return m, SaveCmd(m.Repo, m.Tasks)
		}

	// --- Quit ---
	case key.Matches(msg, km.Quit):
		return m, tea.Quit
	}

	return m, nil
}

// updateInput handles all keypresses in ModeAdd and ModeEdit.
// Navigation keys (j/k) are intentionally not handled here — the text input
// component consumes all key events, preventing keybinding collisions.
func updateInput(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	km := model.DefaultKeyMap()

	switch {

	// --- Confirm ---
	case key.Matches(msg, km.Confirm):
		title := m.Input.Value()
		if title == "" {
			// Empty input — treat as cancel.
			m.Input.Blur()
			m.Mode = model.ModeList
			return m, nil
		}

		if m.Mode == model.ModeAdd {
			m.Tasks = append(m.Tasks, task.New(title))
			m.Cursor = m.Cursor.Bottom(len(m.Tasks))
		} else {
			// ModeEdit — record snapshot before mutation for undo.
			idx := m.Cursor.Index
			m.History.Push(history.Entry{
				Kind:     history.ActionEdit,
				Snapshot: m.Tasks[idx],
				Index:    idx,
			})
			m.Tasks[idx].Title = title
		}

		m.Input.Blur()
		m.Mode = model.ModeList
		return m, SaveCmd(m.Repo, m.Tasks)

	// --- Cancel ---
	case key.Matches(msg, km.Cancel):
		m.Input.Blur()
		m.Mode = model.ModeList
		return m, nil

	// --- All other keys — delegate to the text input component ---
	default:
		var cmd tea.Cmd
		m.Input, cmd = m.Input.Update(msg)
		return m, cmd
	}
}

// --- Undo / Redo helpers ---

// applyUndo pops the undo stack and reverses the recorded mutation.
func applyUndo(m model.Model) model.Model {
	e, ok := m.History.Undo()
	if !ok {
		return m
	}

	switch e.Kind {
	case history.ActionDelete:
		// Re-insert the deleted task at its original index.
		m.Tasks = insertTask(m.Tasks, e.Index, e.Snapshot)
		m.Cursor = m.Cursor.Clamp(len(m.Tasks))

	case history.ActionEdit:
		// Restore the task's previous title and description.
		if e.Index < len(m.Tasks) {
			m.Tasks[e.Index].Title = e.Snapshot.Title
			m.Tasks[e.Index].Description = e.Snapshot.Description
		}
	}

	return m
}

// applyRedo pops the redo stack and re-applies the recorded mutation.
func applyRedo(m model.Model) model.Model {
	e, ok := m.History.Redo()
	if !ok {
		return m
	}

	switch e.Kind {
	case history.ActionDelete:
		// Re-delete the task at the recorded index.
		if e.Index < len(m.Tasks) {
			m.Tasks = removeTask(m.Tasks, e.Index)
			m.Cursor = m.Cursor.Clamp(len(m.Tasks))
		}

	case history.ActionEdit:
		// There is no "after" snapshot stored, so redo for edit is a no-op.
		// A future enhancement could store both before/after snapshots.
	}

	return m
}

// --- Slice helpers ---

// removeTask removes the element at index i from tasks, returning the new slice.
func removeTask(tasks []task.Task, i int) []task.Task {
	result := make([]task.Task, 0, len(tasks)-1)
	result = append(result, tasks[:i]...)
	result = append(result, tasks[i+1:]...)
	return result
}

// insertTask inserts t at index i in tasks, shifting existing elements right.
func insertTask(tasks []task.Task, i int, t task.Task) []task.Task {
	result := make([]task.Task, len(tasks)+1)
	copy(result, tasks[:i])
	result[i] = t
	copy(result[i+1:], tasks[i:])
	return result
}
