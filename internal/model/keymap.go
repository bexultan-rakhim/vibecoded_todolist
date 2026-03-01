package model

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all keybindings for the application.
// Grouping them in a struct makes it straightforward to render the footer
// dynamically — the View layer iterates over the relevant bindings for the
// current mode rather than hardcoding hint strings.
type KeyMap struct {
	// --- List mode ---
	Up     key.Binding
	Down   key.Binding
	Top    key.Binding
	Bottom key.Binding
	Add    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Toggle key.Binding
	Quit   key.Binding

	// --- Add / Edit mode ---
	Confirm key.Binding
	Cancel  key.Binding
}

// DefaultKeyMap returns the standard vim-first keybindings described in the
// design document. Callers can override individual bindings if needed.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j", "down"),
		),
		Top: key.NewBinding(
			key.WithKeys("g"),
			key.WithHelp("g", "top"),
		),
		Bottom: key.NewBinding(
			key.WithKeys("G"),
			key.WithHelp("G", "bottom"),
		),
		Add: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x", "d"),
			key.WithHelp("x", "delete"),
		),
		Toggle: key.NewBinding(
			key.WithKeys("enter", " "),
			key.WithHelp("enter", "toggle"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}
