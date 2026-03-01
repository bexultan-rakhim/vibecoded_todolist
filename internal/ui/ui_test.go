package ui

import (
	"strings"
	"testing"

	"todo/internal/model"
	"todo/internal/navigation"
	"todo/internal/repository"
	"todo/internal/task"
)

// newTestModel creates a minimal Model for UI tests.
func newTestModel(t *testing.T) (model.Model, *repository.InMemoryRepository) {
	t.Helper()
	repo := repository.NewInMemoryRepository()
	m, err := model.New(repo)
	if err != nil {
		t.Fatalf("model.New: %v", err)
	}
	return m, repo
}

// --- truncate ---

func TestTruncate_ShortString_Unchanged(t *testing.T) {
	got := truncate("hello", 10)
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestTruncate_ExactLength_Unchanged(t *testing.T) {
	got := truncate("hello", 5)
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestTruncate_LongString_AppendEllipsis(t *testing.T) {
	got := truncate("hello world", 8)
	if !strings.HasSuffix(got, "…") {
		t.Errorf("expected ellipsis suffix, got %q", got)
	}
	if len([]rune(got)) != 8 {
		t.Errorf("expected length 8, got %d", len([]rune(got)))
	}
}

func TestTruncate_ZeroWidth_ReturnsEmpty(t *testing.T) {
	got := truncate("hello", 0)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestTruncate_WidthOne_ReturnsEllipsis(t *testing.T) {
	got := truncate("hello", 1)
	if got != "…" {
		t.Errorf("expected %q, got %q", "…", got)
	}
}

func TestTruncate_EmptyString_ReturnsEmpty(t *testing.T) {
	got := truncate("", 10)
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

// --- viewportBounds ---

func TestViewportBounds_AllItemsFit(t *testing.T) {
	start, end := viewportBounds(0, 5, 10)
	if start != 0 || end != 5 {
		t.Errorf("expected [0,5), got [%d,%d)", start, end)
	}
}

func TestViewportBounds_CursorAtTop(t *testing.T) {
	start, end := viewportBounds(0, 20, 5)
	if start != 0 || end != 5 {
		t.Errorf("expected [0,5), got [%d,%d)", start, end)
	}
}

func TestViewportBounds_CursorAtBottom(t *testing.T) {
	start, end := viewportBounds(19, 20, 5)
	if start != 15 || end != 20 {
		t.Errorf("expected [15,20), got [%d,%d)", start, end)
	}
}

func TestViewportBounds_CursorInMiddle(t *testing.T) {
	start, end := viewportBounds(10, 20, 5)
	if end-start != 5 {
		t.Errorf("expected window size 5, got %d", end-start)
	}
	if start > 10 || end <= 10 {
		t.Errorf("cursor index 10 must be within [%d,%d)", start, end)
	}
}

func TestViewportBounds_CursorAlwaysInWindow(t *testing.T) {
	total := 30
	height := 7
	for cursor := 0; cursor < total; cursor++ {
		start, end := viewportBounds(cursor, total, height)
		if cursor < start || cursor >= end {
			t.Errorf("cursor %d not in window [%d,%d)", cursor, start, end)
		}
	}
}

func TestViewportBounds_WindowNeverExceedsTotal(t *testing.T) {
	start, end := viewportBounds(2, 4, 10)
	if end > 4 {
		t.Errorf("end %d should not exceed total 4", end)
	}
	if start < 0 {
		t.Errorf("start %d should not be negative", start)
	}
}

// --- RenderList (smoke tests — verify no panic and basic structure) ---

func TestRenderList_EmptyTasks_ContainsHint(t *testing.T) {
	result := RenderList([]task.Task{}, navigation.New(), 80, 20)
	if !strings.Contains(result, "No tasks yet") {
		t.Errorf("expected empty state hint in output, got:\n%s", result)
	}
}

func TestRenderList_WithTasks_ContainsTitles(t *testing.T) {
	tasks := []task.Task{
		task.New("first task"),
		task.New("second task"),
	}
	result := RenderList(tasks, navigation.New(), 80, 20)
	if !strings.Contains(result, "first task") {
		t.Errorf("expected 'first task' in output, got:\n%s", result)
	}
	if !strings.Contains(result, "second task") {
		t.Errorf("expected 'second task' in output, got:\n%s", result)
	}
}

func TestRenderList_FocusedRowContainsCursorSymbol(t *testing.T) {
	tasks := []task.Task{task.New("focused")}
	cursor := navigation.Cursor{Index: 0}
	result := RenderList(tasks, cursor, 80, 20)
	// The cursor symbol glyph should appear in the rendered output.
	if !strings.Contains(result, "󰁔") {
		t.Errorf("expected cursor symbol 󰁔 in focused row, got:\n%s", result)
	}
}

func TestRenderList_ZeroDimensions_NoPanic(t *testing.T) {
	// Should not panic even with degenerate dimensions.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RenderList panicked with zero dimensions: %v", r)
		}
	}()
	RenderList([]task.Task{}, navigation.New(), 0, 0)
}

// --- RenderHeader (smoke tests) ---

func TestRenderHeader_ContainsTitle(t *testing.T) {
	result := RenderHeader([]task.Task{}, 80)
	if !strings.Contains(result, "godo") {
		t.Errorf("expected 'godo' in header, got:\n%s", result)
	}
}

func TestRenderHeader_ShowsCorrectCounts(t *testing.T) {
	tasks := []task.Task{task.New("a"), task.New("b")}
	tasks[0].Complete()

	result := RenderHeader(tasks, 80)
	if !strings.Contains(result, "1 done") {
		t.Errorf("expected '1 done' in header, got:\n%s", result)
	}
	if !strings.Contains(result, "1 todo") {
		t.Errorf("expected '1 todo' in header, got:\n%s", result)
	}
}

// --- RenderModal ---

func TestRenderModal_ListMode_ReturnsBgUnchanged(t *testing.T) {
	m, _ := newTestModel(t)
	m.Mode = model.ModeList
	bg := "background content"

	result := RenderModal(m, bg)
	if result != bg {
		t.Errorf("expected bg content unchanged in ModeList, got:\n%s", result)
	}
}

func TestRenderModal_AddMode_ContainsNewTaskTitle(t *testing.T) {
	m, _ := newTestModel(t)
	m.Mode = model.ModeAdd
	m.Width = 120
	m.Height = 40

	result := RenderModal(m, "bg")
	if !strings.Contains(result, "New Task") {
		t.Errorf("expected 'New Task' in modal, got:\n%s", result)
	}
}

func TestRenderModal_EditMode_ContainsEditTitle(t *testing.T) {
	m, _ := newTestModel(t)
	m.Mode = model.ModeEdit
	m.Width = 120
	m.Height = 40

	result := RenderModal(m, "bg")
	if !strings.Contains(result, "Edit Task") {
		t.Errorf("expected 'Edit Task' in modal, got:\n%s", result)
	}
}

func TestRenderModal_ContainsKeyHints(t *testing.T) {
	m, _ := newTestModel(t)
	m.Mode = model.ModeAdd
	m.Width = 120
	m.Height = 40

	result := RenderModal(m, "bg")
	if !strings.Contains(result, "esc to cancel") {
		t.Errorf("expected key hints in modal, got:\n%s", result)
	}
}
