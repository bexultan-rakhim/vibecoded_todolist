package history

import (
	"testing"

	"todo/internal/task"
)

// helpers

func deleteEntry(index int) Entry {
	return Entry{
		Kind:     ActionDelete,
		Snapshot: task.New("deleted task"),
		Index:    index,
	}
}

func editEntry(index int) Entry {
	return Entry{
		Kind:     ActionEdit,
		Snapshot: task.New("old title"),
		Index:    index,
	}
}

// --- New ---

func TestNew_EmptyStacks(t *testing.T) {
	s := New()
	if s.UndoDepth() != 0 {
		t.Errorf("expected undo depth 0, got %d", s.UndoDepth())
	}
	if s.RedoDepth() != 0 {
		t.Errorf("expected redo depth 0, got %d", s.RedoDepth())
	}
}

func TestNew_CannotUndoOrRedo(t *testing.T) {
	s := New()
	if s.CanUndo() {
		t.Error("expected CanUndo() false on empty stack")
	}
	if s.CanRedo() {
		t.Error("expected CanRedo() false on empty stack")
	}
}

// --- Push ---

func TestPush_IncreasesUndoDepth(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	if s.UndoDepth() != 1 {
		t.Errorf("expected undo depth 1, got %d", s.UndoDepth())
	}
}

func TestPush_ClearsRedoStack(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Undo() // moves entry to redo
	if s.RedoDepth() != 1 {
		t.Fatalf("expected redo depth 1 before second push")
	}
	s.Push(editEntry(1)) // new action should wipe redo
	if s.RedoDepth() != 0 {
		t.Errorf("expected redo stack to be cleared after Push, got depth %d", s.RedoDepth())
	}
}

func TestPush_EvictsOldestWhenFull(t *testing.T) {
	s := New()
	// Push maxDepth+1 entries; the first one should be evicted.
	first := Entry{Kind: ActionDelete, Snapshot: task.New("evicted"), Index: 999}
	s.Push(first)
	for i := 1; i <= maxDepth; i++ {
		s.Push(deleteEntry(i))
	}
	if s.UndoDepth() != maxDepth {
		t.Errorf("expected undo depth %d after overflow, got %d", maxDepth, s.UndoDepth())
	}
	// Drain the stack; the evicted entry (index 999) must not appear.
	for s.CanUndo() {
		e, _ := s.Undo()
		if e.Index == 999 {
			t.Error("oldest entry should have been evicted but was found on the stack")
		}
	}
}

// --- Undo ---

func TestUndo_ReturnsFalseOnEmptyStack(t *testing.T) {
	s := New()
	_, ok := s.Undo()
	if ok {
		t.Error("expected Undo() to return false on empty stack")
	}
}

func TestUndo_ReturnsLastPushedEntry(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Push(editEntry(2))

	e, ok := s.Undo()
	if !ok {
		t.Fatal("expected Undo() to succeed")
	}
	if e.Kind != ActionEdit || e.Index != 2 {
		t.Errorf("expected last pushed entry (edit@2), got kind=%v index=%d", e.Kind, e.Index)
	}
}

func TestUndo_DecreasesUndoDepth(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Push(deleteEntry(1))
	s.Undo()
	if s.UndoDepth() != 1 {
		t.Errorf("expected undo depth 1, got %d", s.UndoDepth())
	}
}

func TestUndo_IncreasesRedoDepth(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Undo()
	if s.RedoDepth() != 1 {
		t.Errorf("expected redo depth 1, got %d", s.RedoDepth())
	}
}

func TestUndo_LIFO_Order(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Push(editEntry(1))
	s.Push(deleteEntry(2))

	for _, wantIndex := range []int{2, 1, 0} {
		e, ok := s.Undo()
		if !ok {
			t.Fatalf("expected Undo() to succeed for index %d", wantIndex)
		}
		if e.Index != wantIndex {
			t.Errorf("expected index %d, got %d", wantIndex, e.Index)
		}
	}
}

// --- Redo ---

func TestRedo_ReturnsFalseOnEmptyStack(t *testing.T) {
	s := New()
	_, ok := s.Redo()
	if ok {
		t.Error("expected Redo() to return false on empty redo stack")
	}
}

func TestRedo_ReturnsUndoneEntry(t *testing.T) {
	s := New()
	s.Push(deleteEntry(3))
	e1, _ := s.Undo()

	e2, ok := s.Redo()
	if !ok {
		t.Fatal("expected Redo() to succeed")
	}
	if e2.Kind != e1.Kind || e2.Index != e1.Index {
		t.Errorf("redo entry differs from undone entry: got %+v, want %+v", e2, e1)
	}
}

func TestRedo_MovesEntryBackToUndo(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Undo()
	s.Redo()
	if s.UndoDepth() != 1 {
		t.Errorf("expected undo depth 1 after redo, got %d", s.UndoDepth())
	}
	if s.RedoDepth() != 0 {
		t.Errorf("expected redo depth 0 after redo, got %d", s.RedoDepth())
	}
}

func TestRedo_LIFO_Order(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Push(editEntry(1))
	s.Push(deleteEntry(2))

	s.Undo() // redo stack: [delete@2]
	s.Undo() // redo stack: [delete@2, edit@1]
	s.Undo() // redo stack: [delete@2, edit@1, delete@0]

	for _, wantIndex := range []int{0, 1, 2} {
		e, ok := s.Redo()
		if !ok {
			t.Fatalf("expected Redo() to succeed for index %d", wantIndex)
		}
		if e.Index != wantIndex {
			t.Errorf("expected index %d, got %d", wantIndex, e.Index)
		}
	}
}

// --- CanUndo / CanRedo ---

func TestCanUndo_TrueAfterPush(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	if !s.CanUndo() {
		t.Error("expected CanUndo() true after Push")
	}
}

func TestCanRedo_TrueAfterUndo(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Undo()
	if !s.CanRedo() {
		t.Error("expected CanRedo() true after Undo")
	}
}

func TestCanRedo_FalseAfterNewPush(t *testing.T) {
	s := New()
	s.Push(deleteEntry(0))
	s.Undo()
	s.Push(editEntry(1)) // wipes redo
	if s.CanRedo() {
		t.Error("expected CanRedo() false after new Push wipes redo stack")
	}
}

// --- ActionKind.String ---

func TestActionKindString(t *testing.T) {
	cases := []struct {
		kind ActionKind
		want string
	}{
		{ActionDelete, "delete"},
		{ActionEdit, "edit"},
		{ActionKind(99), "unknown"},
	}
	for _, c := range cases {
		if got := c.kind.String(); got != c.want {
			t.Errorf("ActionKind(%d).String() = %q, want %q", c.kind, got, c.want)
		}
	}
}

// --- Full undo/redo round-trip ---

func TestRoundTrip_UndoThenRedo(t *testing.T) {
	s := New()
	entries := []Entry{deleteEntry(0), editEntry(1), deleteEntry(2)}
	for _, e := range entries {
		s.Push(e)
	}

	// Undo all three.
	for i := 0; i < 3; i++ {
		if _, ok := s.Undo(); !ok {
			t.Fatalf("Undo %d failed unexpectedly", i)
		}
	}
	if s.UndoDepth() != 0 || s.RedoDepth() != 3 {
		t.Errorf("after 3 undos: undo=%d redo=%d, want undo=0 redo=3", s.UndoDepth(), s.RedoDepth())
	}

	// Redo all three.
	for i := 0; i < 3; i++ {
		if _, ok := s.Redo(); !ok {
			t.Fatalf("Redo %d failed unexpectedly", i)
		}
	}
	if s.UndoDepth() != 3 || s.RedoDepth() != 0 {
		t.Errorf("after 3 redos: undo=%d redo=%d, want undo=3 redo=0", s.UndoDepth(), s.RedoDepth())
	}
}
