package repository

import "todo/internal/task"

// Compile-time assertions: both implementations must satisfy Repository.
// If either type drifts out of compliance, the build fails immediately with a
// clear error — no test run required to catch the regression.
var _ Repository = (*JSONRepository)(nil)
var _ Repository = (*InMemoryRepository)(nil)

// InMemoryRepository is a Repository implementation that stores tasks in a
// plain Go slice. It is not safe for concurrent use and is intended solely
// for use in tests of packages that depend on the Repository interface
// (e.g. the model and update layers), allowing them to avoid any filesystem I/O.
type InMemoryRepository struct {
	tasks []task.Task
}

// NewInMemoryRepository returns an empty InMemoryRepository.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{tasks: []task.Task{}}
}

// Load returns a shallow copy of the stored tasks so callers cannot mutate
// the repository's internal slice directly.
func (r *InMemoryRepository) Load() ([]task.Task, error) {
	cp := make([]task.Task, len(r.tasks))
	copy(cp, r.tasks)
	return cp, nil
}

// Save replaces the stored tasks with a shallow copy of the provided slice.
func (r *InMemoryRepository) Save(tasks []task.Task) error {
	cp := make([]task.Task, len(tasks))
	copy(cp, tasks)
	r.tasks = cp
	return nil
}
