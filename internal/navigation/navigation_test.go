package navigation

import "testing"

// --- New ---

func TestNew_IndexIsZero(t *testing.T) {
	c := New()
	if c.Index != 0 {
		t.Errorf("expected index 0, got %d", c.Index)
	}
}

// --- Down ---

func TestDown_MovesForward(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Down(3)
	if c.Index != 1 {
		t.Errorf("expected index 1, got %d", c.Index)
	}
}

func TestDown_WrapsToTop(t *testing.T) {
	c := Cursor{Index: 2}
	c = c.Down(3) // last item in a 3-item list → should wrap to 0
	if c.Index != 0 {
		t.Errorf("expected index 0 after wrap, got %d", c.Index)
	}
}

func TestDown_EmptyListIsNoop(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Down(0)
	if c.Index != 0 {
		t.Errorf("expected index 0 for empty list, got %d", c.Index)
	}
}

func TestDown_SingleItem_StaysAtZero(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Down(1)
	if c.Index != 0 {
		t.Errorf("expected index 0 for single-item list, got %d", c.Index)
	}
}

// --- Up ---

func TestUp_MovesBackward(t *testing.T) {
	c := Cursor{Index: 2}
	c = c.Up(3)
	if c.Index != 1 {
		t.Errorf("expected index 1, got %d", c.Index)
	}
}

func TestUp_WrapsToBottom(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Up(3) // first item → should wrap to last (index 2)
	if c.Index != 2 {
		t.Errorf("expected index 2 after wrap, got %d", c.Index)
	}
}

func TestUp_EmptyListIsNoop(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Up(0)
	if c.Index != 0 {
		t.Errorf("expected index 0 for empty list, got %d", c.Index)
	}
}

func TestUp_SingleItem_StaysAtZero(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Up(1)
	if c.Index != 0 {
		t.Errorf("expected index 0 for single-item list, got %d", c.Index)
	}
}

// --- Top ---

func TestTop_JumpsToFirst(t *testing.T) {
	c := Cursor{Index: 4}
	c = c.Top(5)
	if c.Index != 0 {
		t.Errorf("expected index 0, got %d", c.Index)
	}
}

func TestTop_EmptyListIsNoop(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Top(0)
	if c.Index != 0 {
		t.Errorf("expected index 0 for empty list, got %d", c.Index)
	}
}

func TestTop_AlreadyAtTop(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Top(5)
	if c.Index != 0 {
		t.Errorf("expected index to remain 0, got %d", c.Index)
	}
}

// --- Bottom ---

func TestBottom_JumpsToLast(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Bottom(5)
	if c.Index != 4 {
		t.Errorf("expected index 4, got %d", c.Index)
	}
}

func TestBottom_EmptyListIsNoop(t *testing.T) {
	c := Cursor{Index: 0}
	c = c.Bottom(0)
	if c.Index != 0 {
		t.Errorf("expected index 0 for empty list, got %d", c.Index)
	}
}

func TestBottom_AlreadyAtBottom(t *testing.T) {
	c := Cursor{Index: 4}
	c = c.Bottom(5)
	if c.Index != 4 {
		t.Errorf("expected index to remain 4, got %d", c.Index)
	}
}

// --- Clamp ---

func TestClamp_IndexInRange_Unchanged(t *testing.T) {
	c := Cursor{Index: 2}
	c = c.Clamp(5)
	if c.Index != 2 {
		t.Errorf("expected index 2, got %d", c.Index)
	}
}

func TestClamp_IndexPastEnd_MovesToLast(t *testing.T) {
	c := Cursor{Index: 5}
	c = c.Clamp(3) // list shrunk to 3 items, valid indices 0-2
	if c.Index != 2 {
		t.Errorf("expected index 2, got %d", c.Index)
	}
}

func TestClamp_EmptyList_ReturnsZero(t *testing.T) {
	c := Cursor{Index: 3}
	c = c.Clamp(0)
	if c.Index != 0 {
		t.Errorf("expected index 0 for empty list, got %d", c.Index)
	}
}

func TestClamp_ExactlyAtLastIndex_Unchanged(t *testing.T) {
	c := Cursor{Index: 4}
	c = c.Clamp(5) // last valid index is 4
	if c.Index != 4 {
		t.Errorf("expected index 4, got %d", c.Index)
	}
}

// --- IsSelected ---

func TestIsSelected_MatchingIndex(t *testing.T) {
	c := Cursor{Index: 2}
	if !c.IsSelected(2) {
		t.Error("expected IsSelected(2) to be true")
	}
}

func TestIsSelected_NonMatchingIndex(t *testing.T) {
	c := Cursor{Index: 2}
	if c.IsSelected(0) || c.IsSelected(1) || c.IsSelected(3) {
		t.Error("expected IsSelected to be false for non-matching indices")
	}
}

// --- Wrapping round-trips ---

func TestDown_FullCircle(t *testing.T) {
	c := New()
	listLen := 5
	for i := 0; i < listLen; i++ {
		c = c.Down(listLen)
	}
	// After exactly listLen steps we should be back at 0.
	if c.Index != 0 {
		t.Errorf("expected full Down circle to return to 0, got %d", c.Index)
	}
}

func TestUp_FullCircle(t *testing.T) {
	c := New()
	listLen := 5
	for i := 0; i < listLen; i++ {
		c = c.Up(listLen)
	}
	// After exactly listLen steps upward we should be back at 0.
	if c.Index != 0 {
		t.Errorf("expected full Up circle to return to 0, got %d", c.Index)
	}
}
