package update

import (
	"errors"
	"testing"

	"todo/internal/repository"
	"todo/internal/task"
)

func TestSaveCmd_SuccessEmitsSavedMsg(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	tasks := []task.Task{task.New("test task")}

	cmd := SaveCmd(repo, tasks)
	msg := cmd() // execute the command synchronously in tests

	if _, ok := msg.(SavedMsg); !ok {
		t.Errorf("expected SavedMsg, got %T: %v", msg, msg)
	}
}

func TestSaveCmd_FailureEmitsSaveErrMsg(t *testing.T) {
	repo := &failingRepository{}
	tasks := []task.Task{task.New("test task")}

	cmd := SaveCmd(repo, tasks)
	msg := cmd()

	errMsg, ok := msg.(SaveErrMsg)
	if !ok {
		t.Fatalf("expected SaveErrMsg, got %T: %v", msg, msg)
	}
	if errMsg.Err == nil {
		t.Error("expected non-nil error in SaveErrMsg")
	}
}

func TestSaveCmd_EmptyTaskList_Succeeds(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	cmd := SaveCmd(repo, []task.Task{})
	msg := cmd()

	if _, ok := msg.(SavedMsg); !ok {
		t.Errorf("expected SavedMsg for empty task list, got %T", msg)
	}
}

func TestSaveCmd_PersistsTasks(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	tasks := []task.Task{task.New("persisted task")}

	cmd := SaveCmd(repo, tasks)
	cmd() // execute

	loaded, err := repo.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Title != "persisted task" {
		t.Errorf("expected persisted task in repo, got %+v", loaded)
	}
}

// --- test doubles ---

// failingRepository is a Repository that always returns an error on Save.
type failingRepository struct{}

func (f *failingRepository) Load() ([]task.Task, error) {
	return []task.Task{}, nil
}

func (f *failingRepository) Save(_ []task.Task) error {
	return errors.New("disk full")
}
