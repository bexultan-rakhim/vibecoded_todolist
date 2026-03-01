package model

import (
	"testing"

	"todo/internal/repository"
	"todo/internal/task"
)

// --- ViewMode.String() ---

func TestViewModeString(t *testing.T) {
	cases := []struct {
		mode ViewMode
		want string
	}{
		{ModeList, "list"},
		{ModeAdd, "add"},
		{ModeEdit, "edit"},
		{ViewMode(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.mode.String(); got != c.want {
			t.Errorf("ViewMode(%d).String() = %q, want %q", c.mode, got, c.want)
		}
	}
}

// --- New() ---

func TestNew_DefaultModeIsList(t *testing.T) {
	m, err := New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if m.Mode != ModeList {
		t.Errorf("expected ModeList, got %v", m.Mode)
	}
}

func TestNew_CursorStartsAtZero(t *testing.T) {
	m, err := New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if m.Cursor.Index != 0 {
		t.Errorf("expected cursor index 0, got %d", m.Cursor.Index)
	}
}

func TestNew_LoadsTasksFromRepository(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	tasks := []task.Task{
		task.New("task one"),
		task.New("task two"),
	}
	if err := repo.Save(tasks); err != nil {
		t.Fatalf("Save: %v", err)
	}

	m, err := New(repo)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	if len(m.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(m.Tasks))
	}
}

func TestNew_EmptyRepositoryGivesEmptyTaskList(t *testing.T) {
	m, err := New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(m.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(m.Tasks))
	}
}

func TestNew_ErrIsNil(t *testing.T) {
	m, err := New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if m.Err != nil {
		t.Errorf("expected Err to be nil, got %v", m.Err)
	}
}

func TestNew_InputHasPlaceholder(t *testing.T) {
	m, err := New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if m.Input.Placeholder == "" {
		t.Error("expected Input to have a placeholder set")
	}
}

func TestNew_RepoIsStored(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	m, err := New(repo)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if m.Repo == nil {
		t.Error("expected Repo to be stored on the model, got nil")
	}
}

// --- Init() ---

func TestInit_ReturnsNilCmd(t *testing.T) {
	m, err := New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil cmd")
	}
}

// --- Mode transitions (state only, no Update logic) ---

func TestModeEnum_ValuesAreDistinct(t *testing.T) {
	modes := []ViewMode{ModeList, ModeAdd, ModeEdit}
	seen := make(map[ViewMode]bool)
	for _, mode := range modes {
		if seen[mode] {
			t.Errorf("duplicate ViewMode value: %d", mode)
		}
		seen[mode] = true
	}
}
