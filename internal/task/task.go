// Package task defines the core domain model for Godo.
// It has no dependencies on storage, UI, or any other internal package —
// it is the innermost layer of the architecture.
package task

import (
	"time"
	"github.com/google/uuid"
)

// Status represents the lifecycle state of a Task.
type Status int

const (
	// StatusTodo is the default state for a newly created task.
	StatusTodo Status = iota
	// StatusDone marks a task as completed.
	StatusDone
)

// String returns a human-readable label for the status.
func (s Status) String() string {
	switch s {
	case StatusTodo:
		return "todo"
	case StatusDone:
		return "done"
	default:
		return "unknown"
	}
}

// Icon returns the Nerd Font icon associated with the status.
// 󰄱 = unchecked box (Todo), 󰄵 = checked box (Done).
func (s Status) Icon() string {
	switch s {
	case StatusTodo:
		return "󰄱"
	case StatusDone:
		return "󰄵"
	default:
		return "?"
	}
}

// Task is the central domain object. All fields are exported so they can be
// serialised to JSON by the repository layer without requiring custom marshalers.
type Task struct {
	// ID is a universally unique identifier assigned at creation time.
	// Using UUID v4 (random) avoids any coordination overhead.
	ID string `json:"id"`

	// Title is the short, one-line summary shown in the list view.
	Title string `json:"title"`

	// Description is an optional, multi-line body shown in the detail/edit modal.
	Description string `json:"description,omitempty"`

	// Status tracks whether the task is pending or complete.
	Status Status `json:"status"`

	// CreatedAt is set once at creation and never mutated.
	CreatedAt time.Time `json:"created_at"`

	// CompletedAt is nil until the task is toggled to StatusDone.
	// A pointer makes absence explicit in JSON (null vs. a zero date).
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// New constructs a Task with a generated UUID and the current timestamp.
// Title is required; Description can be set afterwards via the Edit flow.
func New(title string) Task {
	return Task{
		ID:        uuid.New().String(),
		Title:     title,
		Status:    StatusTodo,
		CreatedAt: time.Now(),
	}
}

// Complete transitions the task to StatusDone and records the completion time.
// Calling Complete on an already-done task is a no-op.
func (t *Task) Complete() {
	if t.Status == StatusDone {
		return
	}
	now := time.Now()
	t.Status = StatusDone
	t.CompletedAt = &now
}

// Reopen transitions a completed task back to StatusTodo and clears CompletedAt.
// Calling Reopen on a todo task is a no-op.
func (t *Task) Reopen() {
	if t.Status == StatusTodo {
		return
	}
	t.Status = StatusTodo
	t.CompletedAt = nil
}

// Toggle flips the task between StatusTodo and StatusDone.
// It is a convenience wrapper used by the Enter/Space keybinding.
func (t *Task) Toggle() {
	if t.Status == StatusTodo {
		t.Complete()
	} else {
		t.Reopen()
	}
}

// IsDone is a convenience predicate for use in view rendering logic.
func (t *Task) IsDone() bool {
	return t.Status == StatusDone
}
