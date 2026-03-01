// Package repository defines the persistence boundary for Godo.
// All storage implementations — JSON, SQLite, in-memory, cloud — must satisfy
// the Repository interface. No other package imports a concrete type directly;
// they depend only on this interface, keeping the UI and business logic
// completely decoupled from "how tasks are stored."
package repository

import "todo/internal/task"

// Repository is the contract that any storage backend must fulfil.
// It is intentionally minimal: the app only needs to load the full task list
// on startup and persist it on every mutation.
type Repository interface {
	// Load retrieves all tasks from the underlying store.
	// If no tasks have been saved yet, an empty slice and a nil error are
	// returned — a missing store is not an error condition.
	// A non-nil error indicates a genuine I/O or decoding failure.
	Load() ([]task.Task, error)

	// Save persists the complete task list to the underlying store,
	// replacing whatever was there before.
	// Implementations must guarantee that a failed Save does not corrupt
	// previously saved data (i.e. they should use an atomic write strategy).
	Save(tasks []task.Task) error
}
