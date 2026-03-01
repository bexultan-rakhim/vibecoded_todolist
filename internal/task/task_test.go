package task

import (
	"testing"
	"time"
)

// --- New() ---

func TestNew_SetsTitle(t *testing.T) {
	tk := New("buy groceries")
	if tk.Title != "buy groceries" {
		t.Errorf("expected title %q, got %q", "buy groceries", tk.Title)
	}
}

func TestNew_DefaultStatusIsTodo(t *testing.T) {
	tk := New("task")
	if tk.Status != StatusTodo {
		t.Errorf("expected StatusTodo, got %v", tk.Status)
	}
}

func TestNew_GeneratesNonEmptyID(t *testing.T) {
	tk := New("task")
	if tk.ID == "" {
		t.Error("expected a non-empty UUID, got empty string")
	}
}

func TestNew_GeneratesUniqueIDs(t *testing.T) {
	a := New("task a")
	b := New("task b")
	if a.ID == b.ID {
		t.Errorf("expected unique IDs, both got %q", a.ID)
	}
}

func TestNew_SetsCreatedAt(t *testing.T) {
	before := time.Now()
	tk := New("task")
	after := time.Now()

	if tk.CreatedAt.Before(before) || tk.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v is outside expected range [%v, %v]", tk.CreatedAt, before, after)
	}
}

func TestNew_CompletedAtIsNil(t *testing.T) {
	tk := New("task")
	if tk.CompletedAt != nil {
		t.Errorf("expected CompletedAt to be nil, got %v", tk.CompletedAt)
	}
}

// --- Complete() ---

func TestComplete_SetsStatusDone(t *testing.T) {
	tk := New("task")
	tk.Complete()
	if tk.Status != StatusDone {
		t.Errorf("expected StatusDone, got %v", tk.Status)
	}
}

func TestComplete_SetsCompletedAt(t *testing.T) {
	before := time.Now()
	tk := New("task")
	tk.Complete()
	after := time.Now()

	if tk.CompletedAt == nil {
		t.Fatal("expected CompletedAt to be set, got nil")
	}
	if tk.CompletedAt.Before(before) || tk.CompletedAt.After(after) {
		t.Errorf("CompletedAt %v is outside expected range [%v, %v]", tk.CompletedAt, before, after)
	}
}

func TestComplete_IsIdempotent(t *testing.T) {
	tk := New("task")
	tk.Complete()
	first := *tk.CompletedAt

	tk.Complete() // second call — should be a no-op
	if *tk.CompletedAt != first {
		t.Error("second Complete() call mutated CompletedAt; expected no-op")
	}
}

// --- Reopen() ---

func TestReopen_SetsStatusTodo(t *testing.T) {
	tk := New("task")
	tk.Complete()
	tk.Reopen()
	if tk.Status != StatusTodo {
		t.Errorf("expected StatusTodo after Reopen, got %v", tk.Status)
	}
}

func TestReopen_ClearsCompletedAt(t *testing.T) {
	tk := New("task")
	tk.Complete()
	tk.Reopen()
	if tk.CompletedAt != nil {
		t.Errorf("expected CompletedAt to be nil after Reopen, got %v", tk.CompletedAt)
	}
}

func TestReopen_IsIdempotentOnTodo(t *testing.T) {
	tk := New("task")
	tk.Reopen() // already todo — should be a no-op, no panic
	if tk.Status != StatusTodo {
		t.Errorf("expected StatusTodo to remain, got %v", tk.Status)
	}
}

// --- Toggle() ---

func TestToggle_TodoBecomeDone(t *testing.T) {
	tk := New("task")
	tk.Toggle()
	if tk.Status != StatusDone {
		t.Errorf("expected StatusDone after Toggle, got %v", tk.Status)
	}
}

func TestToggle_DoneBecomeTodo(t *testing.T) {
	tk := New("task")
	tk.Complete()
	tk.Toggle()
	if tk.Status != StatusTodo {
		t.Errorf("expected StatusTodo after second Toggle, got %v", tk.Status)
	}
}

func TestToggle_RoundTrip(t *testing.T) {
	tk := New("task")
	tk.Toggle() // todo → done
	tk.Toggle() // done → todo
	if tk.Status != StatusTodo {
		t.Errorf("expected StatusTodo after round-trip Toggle, got %v", tk.Status)
	}
	if tk.CompletedAt != nil {
		t.Error("expected CompletedAt to be nil after round-trip Toggle")
	}
}

// --- IsDone() ---

func TestIsDone_FalseForNewTask(t *testing.T) {
	tk := New("task")
	if tk.IsDone() {
		t.Error("expected IsDone() to be false for a new task")
	}
}

func TestIsDone_TrueAfterComplete(t *testing.T) {
	tk := New("task")
	tk.Complete()
	if !tk.IsDone() {
		t.Error("expected IsDone() to be true after Complete()")
	}
}

func TestIsDone_FalseAfterReopen(t *testing.T) {
	tk := New("task")
	tk.Complete()
	tk.Reopen()
	if tk.IsDone() {
		t.Error("expected IsDone() to be false after Reopen()")
	}
}

// --- Status.String() ---

func TestStatusString(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{StatusTodo, "todo"},
		{StatusDone, "done"},
		{Status(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.status.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.status, got, c.want)
		}
	}
}

// --- Status.Icon() ---

func TestStatusIcon(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{StatusTodo, "󰄱"},
		{StatusDone, "󰄵"},
		{Status(99), "?"},
	}
	for _, c := range cases {
		if got := c.status.Icon(); got != c.want {
			t.Errorf("Status(%d).Icon() = %q, want %q", c.status, got, c.want)
		}
	}
}
