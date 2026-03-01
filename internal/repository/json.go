// Package repository provides the persistence layer for Godo.
// This file implements atomic "write-ahead" style file saving to prevent
// JSON corruption if the process crashes or the terminal closes mid-write.
package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"todo/internal/task"
)

// dataFile is the name of the JSON file stored inside the config directory.
const dataFile = "data.json"

// JSONRepository persists tasks as JSON to the local filesystem.
// The zero value is not valid; use NewJSONRepository to construct one.
type JSONRepository struct {
	// dir is the directory where data.json lives (e.g. ~/.config/godo/).
	dir string
}

// NewJSONRepository creates a JSONRepository rooted at the given directory.
// The directory is created (including all parents) if it does not exist.
func NewJSONRepository(dir string) (*JSONRepository, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("repository: create config dir %q: %w", dir, err)
	}
	return &JSONRepository{dir: dir}, nil
}

// dataPath returns the absolute path to the JSON data file.
func (r *JSONRepository) dataPath() string {
	return filepath.Join(r.dir, dataFile)
}

// tempPath returns the path for the temporary write-ahead file.
// It lives in the same directory as the target so that the final os.Rename
// is an atomic same-filesystem operation on all POSIX systems.
func (r *JSONRepository) tempPath() string {
	return filepath.Join(r.dir, "."+dataFile+".tmp")
}

// Load reads and deserialises all tasks from disk.
// If the data file does not exist yet, an empty slice is returned — this is
// the expected state on first run and is not treated as an error.
func (r *JSONRepository) Load() ([]task.Task, error) {
	data, err := os.ReadFile(r.dataPath())
	if os.IsNotExist(err) {
		return []task.Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("repository: read %q: %w", r.dataPath(), err)
	}

	var tasks []task.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("repository: parse %q: %w", r.dataPath(), err)
	}

	return tasks, nil
}

// Save serialises tasks to disk using a write-ahead (atomic rename) strategy:
//
//  1. Marshal the data to JSON.
//  2. Write to a hidden .tmp file in the same directory.
//  3. fsync the tmp file so the OS flushes its buffers to disk.
//  4. Rename tmp → data.json (atomic on POSIX; best-effort on Windows).
//
// This ensures data.json is never left in a partially written state.
// If the process dies between steps 2 and 4, the original data.json is untouched.
func (r *JSONRepository) Save(tasks []task.Task) error {
	// Marshal with indentation so the file is human-readable and diff-friendly.
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("repository: marshal tasks: %w", err)
	}

	// Write to the temporary file.
	tmp := r.tempPath()
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("repository: open tmp file %q: %w", tmp, err)
	}

	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return fmt.Errorf("repository: write tmp file %q: %w", tmp, err)
	}

	// fsync before close to guarantee the OS has flushed to disk.
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("repository: sync tmp file %q: %w", tmp, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("repository: close tmp file %q: %w", tmp, err)
	}

	// Atomic rename: on POSIX this is guaranteed to be atomic by the kernel.
	// On Windows, os.Rename is not atomic but is still safer than direct write.
	if err := os.Rename(tmp, r.dataPath()); err != nil {
		return fmt.Errorf("repository: rename %q → %q: %w", tmp, r.dataPath(), err)
	}

	return nil
}
