// Package navigation implements the cursor logic for Godo's list view.
// It is a pure, stateless package — all functions take the current state as
// arguments and return a new value, with no side effects. This makes the
// logic trivial to test and keeps it fully decoupled from the Bubble Tea model.
package navigation

// Cursor holds the index of the currently selected item in the task list.
// The zero value (Cursor{Index: 0}) is valid and represents the first item.
type Cursor struct {
	Index int
}

// New returns a Cursor pointing at the first item.
func New() Cursor {
	return Cursor{Index: 0}
}

// Down moves the cursor one step toward the end of the list.
// If the cursor is already at the last item it wraps around to the first.
// If the list is empty, the cursor is unchanged.
func (c Cursor) Down(listLen int) Cursor {
	if listLen == 0 {
		return c
	}
	return Cursor{Index: (c.Index + 1) % listLen}
}

// Up moves the cursor one step toward the start of the list.
// If the cursor is already at the first item it wraps around to the last.
// If the list is empty, the cursor is unchanged.
func (c Cursor) Up(listLen int) Cursor {
	if listLen == 0 {
		return c
	}
	return Cursor{Index: (c.Index - 1 + listLen) % listLen}
}

// Top moves the cursor to the first item (gg in Vim).
// If the list is empty, the cursor is unchanged.
func (c Cursor) Top(listLen int) Cursor {
	if listLen == 0 {
		return c
	}
	return Cursor{Index: 0}
}

// Bottom moves the cursor to the last item (G in Vim).
// If the list is empty, the cursor is unchanged.
func (c Cursor) Bottom(listLen int) Cursor {
	if listLen == 0 {
		return c
	}
	return Cursor{Index: listLen - 1}
}

// Clamp ensures the cursor index stays within bounds after an external mutation
// (e.g. a task deletion that shortens the list). It always returns a valid index:
//   - Empty list  → 0 (callers should check listLen before rendering)
//   - Index past end → last item
func (c Cursor) Clamp(listLen int) Cursor {
	if listLen == 0 {
		return Cursor{Index: 0}
	}
	if c.Index >= listLen {
		return Cursor{Index: listLen - 1}
	}
	return c
}

// IsSelected reports whether the given list index is the currently focused row.
// Intended for use inside the View render loop:
//
//	for i, t := range tasks {
//	    if cursor.IsSelected(i) { ... render highlighted row ... }
//	}
func (c Cursor) IsSelected(index int) bool {
	return c.Index == index
}
