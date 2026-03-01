// Package ui contains all rendering logic for Godo.
// Components are pure functions: they take state and return styled strings.
// No business logic lives here — only layout and presentation.
package ui

import "github.com/charmbracelet/lipgloss"

// Tokyo Night palette — https://github.com/enkia/tokyo-night-vscode-theme
// All colours are defined as constants so they can be referenced by name
// throughout the UI package rather than scattering hex literals everywhere.
const (
	colorBg        = lipgloss.Color("#1a1b2e") // deep navy background
	colorBgAlt     = lipgloss.Color("#16213e") // slightly lighter surface
	colorBgHighlight = lipgloss.Color("#2a2b3d") // focused row background
	colorBorder    = lipgloss.Color("#414868") // muted border
	colorBorderFocus = lipgloss.Color("#7aa2f7") // bright blue — active border
	colorFg        = lipgloss.Color("#c0caf5") // primary text
	colorFgMuted   = lipgloss.Color("#565f89") // dimmed / secondary text
	colorPurple    = lipgloss.Color("#bb9af7") // accent — cursor indicator
	colorBlue      = lipgloss.Color("#7aa2f7") // links, active elements
	colorCyan      = lipgloss.Color("#7dcfff") // highlights
	colorGreen     = lipgloss.Color("#9ece6a") // done status
	colorRed       = lipgloss.Color("#f7768e") // errors, delete hints
	colorYellow    = lipgloss.Color("#e0af68") // warnings, edit hints
)

// --- Shared base styles ---

// Base is the root style applied to the entire application window.
var Base = lipgloss.NewStyle().
	Background(colorBg).
	Foreground(colorFg)

// BorderBox is the standard rounded border used for the main layout panels.
var BorderBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder)

// FocusedBorderBox is BorderBox with a highlighted border colour,
// used when a panel or modal is the active focus target.
var FocusedBorderBox = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorderFocus)

// --- Header styles ---

var HeaderContainer = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder).
	Padding(0, 1)

var HeaderTitle = lipgloss.NewStyle().
	Foreground(colorPurple).
	Bold(true)

var HeaderStats = lipgloss.NewStyle().
	Foreground(colorFgMuted)

var HeaderStatsDone = lipgloss.NewStyle().
	Foreground(colorGreen)

// --- List / body styles ---

var ListContainer = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder)

// RowNormal is the style for an unselected task row.
var RowNormal = lipgloss.NewStyle().
	Foreground(colorFg).
	Padding(0, 1)

// RowFocused is the style for the currently selected task row.
// Electric purple background with a leading accent colour.
var RowFocused = lipgloss.NewStyle().
	Background(colorBgHighlight).
	Foreground(colorCyan).
	Bold(true).
	Padding(0, 1)

// RowDone is applied to completed tasks in normal (unfocused) rows.
var RowDone = lipgloss.NewStyle().
	Foreground(colorFgMuted).
	Strikethrough(true).
	Padding(0, 1)

// RowDoneFocused is a completed task that is also the selected row.
var RowDoneFocused = lipgloss.NewStyle().
	Background(colorBgHighlight).
	Foreground(colorFgMuted).
	Strikethrough(true).
	Bold(true).
	Padding(0, 1)

// CursorIndicator is the leading symbol rendered beside the focused row.
var CursorIndicator = lipgloss.NewStyle().
	Foreground(colorPurple).
	Bold(true)

// StatusIconTodo styles the 󰄱 icon.
var StatusIconTodo = lipgloss.NewStyle().
	Foreground(colorFgMuted)

// StatusIconDone styles the 󰄵 icon.
var StatusIconDone = lipgloss.NewStyle().
	Foreground(colorGreen)

// EmptyState is used when no tasks exist.
var EmptyState = lipgloss.NewStyle().
	Foreground(colorFgMuted).
	Italic(true).
	Padding(1, 2)

// --- Footer styles ---

var FooterContainer = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorder).
	Padding(0, 1)

// KeyHintKey styles the key portion of a footer hint (e.g. "a").
var KeyHintKey = lipgloss.NewStyle().
	Foreground(colorBlue).
	Bold(true)

// KeyHintDesc styles the description portion of a footer hint (e.g. "add").
var KeyHintDesc = lipgloss.NewStyle().
	Foreground(colorFgMuted)

// FooterSeparator is the divider between key hints.
var FooterSeparator = lipgloss.NewStyle().
	Foreground(colorBorder)

// FooterError renders an error message in the footer area.
var FooterError = lipgloss.NewStyle().
	Foreground(colorRed).
	Bold(true)

// --- Modal styles ---

var ModalContainer = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorBorderFocus).
	Background(colorBgAlt).
	Padding(1, 2)

var ModalTitle = lipgloss.NewStyle().
	Foreground(colorPurple).
	Bold(true).
	MarginBottom(1)

var ModalHint = lipgloss.NewStyle().
	Foreground(colorFgMuted).
	Italic(true)
