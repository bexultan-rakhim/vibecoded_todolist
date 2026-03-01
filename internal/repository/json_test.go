package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"todo/internal/task"
)

// newTestRepo creates a JSONRepository backed by a temporary directory that is
// automatically cleaned up when the test finishes.
func newTestRepo(t *testing.T) *JSONRepository {
	t.Helper()
	dir := t.TempDir()
	repo, err := NewJSONRepository(dir)
	if err != nil {
		t.Fatalf("NewJSONRepository: %v", err)
	}
	return repo
}

// sampleTasks returns a small, deterministic slice of tasks for use in tests.
func sampleTasks() []task.Task {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	done := time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC)

	return []task.Task{
		{
			ID:        "uuid-1",
			Title:     "Buy groceries",
			Status:    task.StatusTodo,
			CreatedAt: now,
		},
		{
			ID:          "uuid-2",
			Title:       "Write tests",
			Description: "Always write tests",
			Status:      task.StatusDone,
			CreatedAt:   now,
			CompletedAt: &done,
		},
	}
}

// --- NewJSONRepository ---

func TestNewJSONRepository_CreatesDirectory(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "nested", "config", "godo")

	_, err := NewJSONRepository(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("expected directory %q to be created, but it does not exist", dir)
	}
}

func TestNewJSONRepository_ExistingDirectoryIsOk(t *testing.T) {
	dir := t.TempDir() // already exists
	_, err := NewJSONRepository(dir)
	if err != nil {
		t.Errorf("expected no error for existing directory, got %v", err)
	}
}

// --- Load ---

func TestLoad_EmptySliceWhenFileAbsent(t *testing.T) {
	repo := newTestRepo(t)

	tasks, err := repo.Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected empty slice, got %d tasks", len(tasks))
	}
}

func TestLoad_ReturnsPersistedTasks(t *testing.T) {
	repo := newTestRepo(t)
	want := sampleTasks()

	if err := repo.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("expected %d tasks, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i].ID != want[i].ID {
			t.Errorf("task[%d].ID: want %q, got %q", i, want[i].ID, got[i].ID)
		}
		if got[i].Title != want[i].Title {
			t.Errorf("task[%d].Title: want %q, got %q", i, want[i].Title, got[i].Title)
		}
		if got[i].Description != want[i].Description {
			t.Errorf("task[%d].Description: want %q, got %q", i, want[i].Description, got[i].Description)
		}
		if got[i].Status != want[i].Status {
			t.Errorf("task[%d].Status: want %v, got %v", i, want[i].Status, got[i].Status)
		}
	}
}

func TestLoad_PreservesCompletedAt(t *testing.T) {
	repo := newTestRepo(t)
	tasks := sampleTasks()

	if err := repo.Save(tasks); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// tasks[0] has no CompletedAt
	if got[0].CompletedAt != nil {
		t.Errorf("expected CompletedAt to be nil for todo task, got %v", got[0].CompletedAt)
	}

	// tasks[1] has a CompletedAt
	if got[1].CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set for done task, got nil")
	}
	if !got[1].CompletedAt.Equal(*tasks[1].CompletedAt) {
		t.Errorf("CompletedAt mismatch: want %v, got %v", tasks[1].CompletedAt, got[1].CompletedAt)
	}
}

func TestLoad_ErrorOnMalformedJSON(t *testing.T) {
	repo := newTestRepo(t)

	// Write deliberately broken JSON directly to the data file.
	if err := os.WriteFile(repo.dataPath(), []byte(`{not valid json`), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	_, err := repo.Load()
	if err == nil {
		t.Error("expected an error for malformed JSON, got nil")
	}
}

// --- Save ---

func TestSave_CreatesDataFile(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Save(sampleTasks()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(repo.dataPath()); os.IsNotExist(err) {
		t.Errorf("expected data file %q to exist after Save", repo.dataPath())
	}
}

func TestSave_EmptySlice(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Save([]task.Task{}); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected 0 tasks after saving empty slice, got %d", len(got))
	}
}

func TestSave_OverwritesPreviousData(t *testing.T) {
	repo := newTestRepo(t)

	first := sampleTasks()
	if err := repo.Save(first); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	second := []task.Task{task.New("only this one")}
	if err := repo.Save(second); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 task after overwrite, got %d", len(got))
	}
	if got[0].Title != "only this one" {
		t.Errorf("unexpected task title %q", got[0].Title)
	}
}

func TestSave_ProducesValidJSON(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Save(sampleTasks()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	raw, err := os.ReadFile(repo.dataPath())
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var parsed []map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		t.Errorf("data file contains invalid JSON: %v", err)
	}
}

func TestSave_NoTempFileLeftBehind(t *testing.T) {
	repo := newTestRepo(t)

	if err := repo.Save(sampleTasks()); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if _, err := os.Stat(repo.tempPath()); !os.IsNotExist(err) {
		t.Errorf("expected tmp file %q to be removed after Save, but it exists", repo.tempPath())
	}
}

// --- Round-trip ---

func TestRoundTrip_SaveThenLoadIsIdentical(t *testing.T) {
	repo := newTestRepo(t)
	want := sampleTasks()

	if err := repo.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("length mismatch: want %d, got %d", len(want), len(got))
	}

	for i := range want {
		w, g := want[i], got[i]
		if w.ID != g.ID || w.Title != g.Title || w.Status != g.Status {
			t.Errorf("task[%d] mismatch:\n  want: %+v\n  got:  %+v", i, w, g)
		}
	}
}
