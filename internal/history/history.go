// Package history implements a bounded undo/redo stack for Godo.
//
// Design decisions:
//   - Only Delete and Edit mutations are tracked, per spec.
//   - The stack stores snapshots of the affected task only (not the full list),
//     keeping memory usage minimal.
//   - Maximum depth is 25 entries; the oldest entry is evicted when exceeded.
//   - State is in-memory only and is lost on quit — intentional by design.
package history

import "todo/internal/task"

const maxDepth = 25

// ActionKind identifies the type of mutation that was recorded.
type ActionKind int

const (
	// ActionDelete records that a task was removed from the list.
	// Undo restores the task at its original index.
	ActionDelete ActionKind = iota

	// ActionEdit records that a task's title or description was changed.
	// Undo restores the task's previous title/description in place.
	ActionEdit
)

// String returns a human-readable label for the action kind.
func (a ActionKind) String() string {
	switch a {
	case ActionDelete:
		return "delete"
	case ActionEdit:
		return "edit"
	default:
		return "unknown"
	}
}

// Entry is a single item on the undo or redo stack.
// It captures everything needed to reverse or reapply the mutation.
type Entry struct {
	// Kind identifies what happened so the update layer knows how to reverse it.
	Kind ActionKind

	// Snapshot is the state of the task before the mutation was applied.
	// For ActionDelete: the full task that was removed.
	// For ActionEdit:   the task with its previous title/description.
	Snapshot task.Task

	// Index is the position in the task list where the mutation occurred.
	// For ActionDelete: the index the task occupied before removal.
	// For ActionEdit:   the index of the task that was edited (unchanged).
	Index int
}

// Stack is a bounded, in-memory undo/redo manager.
// The zero value is not valid; use New() to construct one.
type Stack struct {
	undo []Entry
	redo []Entry
}

// New returns an empty Stack ready for use.
func New() *Stack {
	return &Stack{
		undo: make([]Entry, 0, maxDepth),
		redo: make([]Entry, 0, maxDepth),
	}
}

// Push records a new mutation on the undo stack and clears the redo stack.
// Clearing redo is the standard behaviour: once you take a new action after
// an undo, the previously undone future is discarded.
// If the stack is already at maxDepth, the oldest entry is evicted.
func (s *Stack) Push(e Entry) {
	if len(s.undo) >= maxDepth {
		// Evict the oldest entry (front of the slice).
		s.undo = s.undo[1:]
	}
	s.undo = append(s.undo, e)
	// Any branching redo history is now invalid.
	s.redo = s.redo[:0]
}

// Undo pops the most recent entry from the undo stack, pushes it onto the
// redo stack, and returns it so the update layer can reverse the mutation.
// Returns (Entry{}, false) if there is nothing to undo.
func (s *Stack) Undo() (Entry, bool) {
	if len(s.undo) == 0 {
		return Entry{}, false
	}
	last := len(s.undo) - 1
	e := s.undo[last]
	s.undo = s.undo[:last]

	// Cap redo depth as well to stay consistent.
	if len(s.redo) >= maxDepth {
		s.redo = s.redo[1:]
	}
	s.redo = append(s.redo, e)
	return e, true
}

// Redo pops the most recent entry from the redo stack, pushes it back onto
// the undo stack, and returns it so the update layer can reapply the mutation.
// Returns (Entry{}, false) if there is nothing to redo.
func (s *Stack) Redo() (Entry, bool) {
	if len(s.redo) == 0 {
		return Entry{}, false
	}
	last := len(s.redo) - 1
	e := s.redo[last]
	s.redo = s.redo[:last]

	if len(s.undo) >= maxDepth {
		s.undo = s.undo[1:]
	}
	s.undo = append(s.undo, e)
	return e, true
}

// CanUndo reports whether there is at least one action to undo.
// Useful for conditionally rendering the undo hint in the footer.
func (s *Stack) CanUndo() bool {
	return len(s.undo) > 0
}

// CanRedo reports whether there is at least one action to redo.
func (s *Stack) CanRedo() bool {
	return len(s.redo) > 0
}

// UndoDepth returns the number of entries currently on the undo stack.
func (s *Stack) UndoDepth() int {
	return len(s.undo)
}

// RedoDepth returns the number of entries currently on the redo stack.
func (s *Stack) RedoDepth() int {
	return len(s.redo)
}
