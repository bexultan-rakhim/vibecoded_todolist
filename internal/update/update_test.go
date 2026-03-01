package update

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"todo/internal/model"
	"todo/internal/repository"
	"todo/internal/task"
)

// --- helpers ---

func newModel(t *testing.T) model.Model {
	t.Helper()
	m, err := model.New(repository.NewInMemoryRepository())
	if err != nil {
		t.Fatalf("model.New: %v", err)
	}
	m.Width = 120
	m.Height = 40
	return m
}

func withTasks(t *testing.T, titles ...string) model.Model {
	t.Helper()
	m := newModel(t)
	for _, title := range titles {
		m.Tasks = append(m.Tasks, task.New(title))
	}
	return m
}

func keyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func specialKey(t tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: t}
}

// sendKey is a convenience wrapper that calls Update with a key message.
func sendKey(m model.Model, key string) (model.Model, tea.Cmd) {
	return Update(m, keyMsg(key))
}

// --- WindowSizeMsg ---

func TestUpdate_WindowSizeMsg_UpdatesDimensions(t *testing.T) {
	m := newModel(t)
	m2, _ := Update(m, tea.WindowSizeMsg{Width: 200, Height: 50})
	if m2.Width != 200 || m2.Height != 50 {
		t.Errorf("expected 200x50, got %dx%d", m2.Width, m2.Height)
	}
}

// --- SavedMsg / SaveErrMsg ---

func TestUpdate_SavedMsg_ClearsError(t *testing.T) {
	m := newModel(t)
	m.Err = errStub
	m2, _ := Update(m, SavedMsg{})
	if m2.Err != nil {
		t.Errorf("expected Err to be cleared, got %v", m2.Err)
	}
}

func TestUpdate_SaveErrMsg_SetsError(t *testing.T) {
	m := newModel(t)
	m2, _ := Update(m, SaveErrMsg{Err: errStub})
	if m2.Err == nil {
		t.Error("expected Err to be set after SaveErrMsg")
	}
}

// --- Navigation ---

func TestUpdate_Down_MovesCursor(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m2, _ := sendKey(m, "j")
	if m2.Cursor.Index != 1 {
		t.Errorf("expected cursor 1, got %d", m2.Cursor.Index)
	}
}

func TestUpdate_Up_MovesCursor(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m.Cursor.Index = 2
	m2, _ := sendKey(m, "k")
	if m2.Cursor.Index != 1 {
		t.Errorf("expected cursor 1, got %d", m2.Cursor.Index)
	}
}

func TestUpdate_Down_WrapsToTop(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m.Cursor.Index = 2
	m2, _ := sendKey(m, "j")
	if m2.Cursor.Index != 0 {
		t.Errorf("expected cursor to wrap to 0, got %d", m2.Cursor.Index)
	}
}

func TestUpdate_Up_WrapsToBottom(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m2, _ := sendKey(m, "k")
	if m2.Cursor.Index != 2 {
		t.Errorf("expected cursor to wrap to 2, got %d", m2.Cursor.Index)
	}
}

func TestUpdate_Top_JumpsToCursor0(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m.Cursor.Index = 2
	m2, _ := sendKey(m, "g")
	if m2.Cursor.Index != 0 {
		t.Errorf("expected cursor 0, got %d", m2.Cursor.Index)
	}
}

func TestUpdate_Bottom_JumpsToLastCursor(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m2, _ := sendKey(m, "G")
	if m2.Cursor.Index != 2 {
		t.Errorf("expected cursor 2, got %d", m2.Cursor.Index)
	}
}

// --- Add mode ---

func TestUpdate_Add_SwitchesToModeAdd(t *testing.T) {
	m := newModel(t)
	m2, _ := sendKey(m, "a")
	if m2.Mode != model.ModeAdd {
		t.Errorf("expected ModeAdd, got %v", m2.Mode)
	}
}

func TestUpdate_Add_ClearsInput(t *testing.T) {
	m := newModel(t)
	m.Input.SetValue("leftover")
	m2, _ := sendKey(m, "a")
	if m2.Input.Value() != "" {
		t.Errorf("expected empty input on ModeAdd entry, got %q", m2.Input.Value())
	}
}

func TestUpdate_AddConfirm_CreatesTask(t *testing.T) {
	m := newModel(t)
	m, _ = sendKey(m, "a")        // enter add mode
	m.Input.SetValue("new task")  // simulate typing
	m2, _ := Update(m, specialKey(tea.KeyEnter))

	if len(m2.Tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(m2.Tasks))
	}
	if m2.Tasks[0].Title != "new task" {
		t.Errorf("expected title 'new task', got %q", m2.Tasks[0].Title)
	}
}

func TestUpdate_AddConfirm_ReturnsModeList(t *testing.T) {
	m := newModel(t)
	m, _ = sendKey(m, "a")
	m.Input.SetValue("task")
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if m2.Mode != model.ModeList {
		t.Errorf("expected ModeList after confirm, got %v", m2.Mode)
	}
}

func TestUpdate_AddConfirm_EmptyInput_Cancels(t *testing.T) {
	m := newModel(t)
	m, _ = sendKey(m, "a")
	// Don't set any input value
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if m2.Mode != model.ModeList {
		t.Errorf("expected ModeList on empty confirm, got %v", m2.Mode)
	}
	if len(m2.Tasks) != 0 {
		t.Errorf("expected no tasks on empty confirm, got %d", len(m2.Tasks))
	}
}

func TestUpdate_AddCancel_ReturnsModeList(t *testing.T) {
	m := newModel(t)
	m, _ = sendKey(m, "a")
	m2, _ := Update(m, specialKey(tea.KeyEscape))
	if m2.Mode != model.ModeList {
		t.Errorf("expected ModeList after cancel, got %v", m2.Mode)
	}
}

func TestUpdate_AddCancel_DoesNotCreateTask(t *testing.T) {
	m := newModel(t)
	m, _ = sendKey(m, "a")
	m.Input.SetValue("abandoned")
	m2, _ := Update(m, specialKey(tea.KeyEscape))
	if len(m2.Tasks) != 0 {
		t.Errorf("expected no tasks after cancel, got %d", len(m2.Tasks))
	}
}

// --- Edit mode ---

func TestUpdate_Edit_SwitchesToModeEdit(t *testing.T) {
	m := withTasks(t, "original")
	m2, _ := sendKey(m, "e")
	if m2.Mode != model.ModeEdit {
		t.Errorf("expected ModeEdit, got %v", m2.Mode)
	}
}

func TestUpdate_Edit_PrePopulatesInput(t *testing.T) {
	m := withTasks(t, "original title")
	m2, _ := sendKey(m, "e")
	if m2.Input.Value() != "original title" {
		t.Errorf("expected input pre-populated with 'original title', got %q", m2.Input.Value())
	}
}

func TestUpdate_EditConfirm_UpdatesTitle(t *testing.T) {
	m := withTasks(t, "old title")
	m, _ = sendKey(m, "e")
	m.Input.SetValue("new title")
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if m2.Tasks[0].Title != "new title" {
		t.Errorf("expected 'new title', got %q", m2.Tasks[0].Title)
	}
}

func TestUpdate_EditConfirm_PushesHistoryEntry(t *testing.T) {
	m := withTasks(t, "original")
	m, _ = sendKey(m, "e")
	m.Input.SetValue("modified")
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if !m2.History.CanUndo() {
		t.Error("expected undo to be available after edit")
	}
}

func TestUpdate_Edit_EmptyList_IsNoop(t *testing.T) {
	m := newModel(t)
	m2, _ := sendKey(m, "e")
	if m2.Mode != model.ModeList {
		t.Errorf("expected ModeList when editing empty list, got %v", m2.Mode)
	}
}

// --- Delete ---

func TestUpdate_Delete_RemovesTask(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m2, _ := sendKey(m, "x")
	if len(m2.Tasks) != 2 {
		t.Errorf("expected 2 tasks after delete, got %d", len(m2.Tasks))
	}
}

func TestUpdate_Delete_Clampscursor(t *testing.T) {
	m := withTasks(t, "only one")
	m2, _ := sendKey(m, "x")
	if m2.Cursor.Index != 0 {
		t.Errorf("expected cursor clamped to 0, got %d", m2.Cursor.Index)
	}
}

func TestUpdate_Delete_PushesHistoryEntry(t *testing.T) {
	m := withTasks(t, "task")
	m2, _ := sendKey(m, "x")
	if !m2.History.CanUndo() {
		t.Error("expected undo to be available after delete")
	}
}

func TestUpdate_Delete_EmptyList_IsNoop(t *testing.T) {
	m := newModel(t)
	m2, _ := sendKey(m, "x")
	if len(m2.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(m2.Tasks))
	}
}

// --- Toggle ---

func TestUpdate_Toggle_MarksTaskDone(t *testing.T) {
	m := withTasks(t, "task")
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if !m2.Tasks[0].IsDone() {
		t.Error("expected task to be done after toggle")
	}
}

func TestUpdate_Toggle_MarksTaskTodo(t *testing.T) {
	m := withTasks(t, "task")
	m.Tasks[0].Complete()
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if m2.Tasks[0].IsDone() {
		t.Error("expected task to be todo after second toggle")
	}
}

func TestUpdate_Toggle_EmptyList_IsNoop(t *testing.T) {
	m := newModel(t)
	m2, _ := Update(m, specialKey(tea.KeyEnter))
	if len(m2.Tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(m2.Tasks))
	}
}

// --- Undo ---

func TestUpdate_Undo_RestoresDeletedTask(t *testing.T) {
	m := withTasks(t, "keep me")
	m, _ = sendKey(m, "x")   // delete
	m2, _ := sendKey(m, "u") // undo
	if len(m2.Tasks) != 1 {
		t.Fatalf("expected 1 task after undo, got %d", len(m2.Tasks))
	}
	if m2.Tasks[0].Title != "keep me" {
		t.Errorf("expected 'keep me', got %q", m2.Tasks[0].Title)
	}
}

func TestUpdate_Undo_RestoresEditedTitle(t *testing.T) {
	m := withTasks(t, "original")
	m, _ = sendKey(m, "e")
	m.Input.SetValue("changed")
	m, _ = Update(m, specialKey(tea.KeyEnter)) // confirm edit
	m2, _ := sendKey(m, "u")                  // undo edit
	if m2.Tasks[0].Title != "original" {
		t.Errorf("expected 'original' after undo, got %q", m2.Tasks[0].Title)
	}
}

func TestUpdate_Undo_EmptyStack_IsNoop(t *testing.T) {
	m := withTasks(t, "task")
	m2, _ := sendKey(m, "u")
	if len(m2.Tasks) != 1 {
		t.Errorf("expected tasks unchanged, got %d", len(m2.Tasks))
	}
}

// --- Redo ---

func TestUpdate_Redo_ReappliesDelete(t *testing.T) {
	m := withTasks(t, "task")
	m, _ = sendKey(m, "x")   // delete
	m, _ = sendKey(m, "u")   // undo (restore)
	m2, _ := sendKey(m, "r") // redo (delete again)
	if len(m2.Tasks) != 0 {
		t.Errorf("expected 0 tasks after redo delete, got %d", len(m2.Tasks))
	}
}

func TestUpdate_Redo_EmptyStack_IsNoop(t *testing.T) {
	m := withTasks(t, "task")
	m2, _ := sendKey(m, "r")
	if len(m2.Tasks) != 1 {
		t.Errorf("expected tasks unchanged, got %d", len(m2.Tasks))
	}
}

// --- Quit ---

func TestUpdate_Quit_ReturnsQuitCmd(t *testing.T) {
	m := newModel(t)
	_, cmd := sendKey(m, "q")
	if cmd == nil {
		t.Error("expected non-nil cmd (tea.Quit) after q")
	}
}

// --- Input mode blocks navigation ---

func TestUpdate_InputMode_NavigationKeysDoNotMoveCursor(t *testing.T) {
	m := withTasks(t, "a", "b", "c")
	m, _ = sendKey(m, "a") // enter add mode
	initialCursor := m.Cursor.Index
	m2, _ := sendKey(m, "j") // j should type into input, not navigate
	if m2.Cursor.Index != initialCursor {
		t.Errorf("cursor moved in input mode: was %d, now %d", initialCursor, m2.Cursor.Index)
	}
}

// --- slice helpers ---

func TestRemoveTask_RemovesCorrectIndex(t *testing.T) {
	tasks := []task.Task{task.New("a"), task.New("b"), task.New("c")}
	result := removeTask(tasks, 1)
	if len(result) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(result))
	}
	if result[0].Title != "a" || result[1].Title != "c" {
		t.Errorf("unexpected titles: %q %q", result[0].Title, result[1].Title)
	}
}

func TestInsertTask_InsertsAtCorrectIndex(t *testing.T) {
	tasks := []task.Task{task.New("a"), task.New("c")}
	inserted := task.New("b")
	result := insertTask(tasks, 1, inserted)
	if len(result) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(result))
	}
	if result[1].Title != "b" {
		t.Errorf("expected 'b' at index 1, got %q", result[1].Title)
	}
}

func TestInsertTask_AtStart(t *testing.T) {
	tasks := []task.Task{task.New("b")}
	result := insertTask(tasks, 0, task.New("a"))
	if result[0].Title != "a" {
		t.Errorf("expected 'a' at index 0, got %q", result[0].Title)
	}
}

func TestInsertTask_AtEnd(t *testing.T) {
	tasks := []task.Task{task.New("a")}
	result := insertTask(tasks, 1, task.New("b"))
	if result[1].Title != "b" {
		t.Errorf("expected 'b' at index 1, got %q", result[1].Title)
	}
}

// --- stubs ---

var errStub = errors.New("stub error")
