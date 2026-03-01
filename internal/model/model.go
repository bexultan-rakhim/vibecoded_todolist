// Package model defines the single source of truth for Godo's application state.
// It follows The Elm Architecture: the Model is a plain data structure with no
// methods other than Init(). All state transitions happen in the update package.
package model

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"todo/internal/history"
	"todo/internal/navigation"
	"todo/internal/repository"
	"todo/internal/task"
)

// ViewMode represents the current interaction state of the application.
// Only one mode is active at a time, and each mode owns a distinct set of
// keybindings — preventing collisions between navigation keys and text input.
type ViewMode int

const (
	// ModeList is the default mode. The task list is visible and fully
	// interactive: j/k navigate, a/e/x/dd act on tasks.
	ModeList ViewMode = iota

	// ModeAdd is active when the user presses 'a'. A modal overlay is shown
	// with a focused text input. Navigation keys (j/k) are disabled.
	ModeAdd

	// ModeEdit is active when the user presses 'e' on a selected task.
	// Identical to ModeAdd visually, but the input is pre-populated with the
	// current task title.
	ModeEdit
)

// String returns a human-readable label for the mode, useful for debugging.
func (m ViewMode) String() string {
	switch m {
	case ModeList:
		return "list"
	case ModeAdd:
		return "add"
	case ModeEdit:
		return "edit"
	default:
		return "unknown"
	}
}

// Model is the single source of truth for the entire application.
// Every field that affects rendering must live here — the View function
// is a pure function of this struct.
type Model struct {
	// Tasks is the ordered list of all tasks. The repository layer owns
	// persistence; this slice is the in-memory working copy.
	Tasks []task.Task

	// Cursor tracks which task is currently highlighted in the list.
	Cursor navigation.Cursor

	// Mode determines which keybindings are active and what the UI renders.
	Mode ViewMode

	// Input is the Bubble Tea text input component used in ModeAdd and ModeEdit.
	// It is always present in the model but only focused and rendered in modal modes.
	Input textinput.Model

	// Repo is the persistence backend. Stored on the model so the update
	// layer can trigger saves without needing a global or closure.
	Repo repository.Repository

	// Width and Height are the current terminal dimensions, updated on every
	// tea.WindowSizeMsg. The View function uses them for responsive layout.
	Width  int
	Height int

	// History is the undo/redo stack. Only Delete and Edit mutations are
	// tracked. It is in-memory only and resets on each app launch.
	History *history.Stack

	// Err holds the last error encountered (e.g. a failed save). When non-nil
	// the footer renders an error message instead of the keymap hint.
	Err error
}

// New constructs the initial Model, loading tasks from the repository.
// It returns an error if the initial Load fails, so main.go can exit cleanly
// with a message instead of starting a broken TUI.
func New(repo repository.Repository) (Model, error) {
	tasks, err := repo.Load()
	if err != nil {
		return Model{}, err
	}

	input := textinput.New()
	input.Placeholder = "Task title..."
	input.CharLimit = 120

	return Model{
		Tasks:  tasks,
		Cursor: navigation.New(),
		Mode:   ModeList,
		Input:  input,
		Repo:   repo,
		History: history.New(),
	}, nil
}

// Init is the Bubble Tea lifecycle method called once at startup.
// We return nil (no initial command) because the repository load happens
// synchronously in New() before the program starts.
func (m Model) Init() tea.Cmd {
	return nil
}
