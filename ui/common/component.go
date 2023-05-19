package common

import (
	// "github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Sizable interface {
	SetSize(width, height int)
	// Width() int
	// Height() int
	// GetMargins() int
}

type Component interface {
	tea.Model
	Sizable
	// help.KeyMap
}

// type Component interface {
// 	tea.Model
// 	Sizable
// }

type FocusableComponent interface {
	Component
	Focused() bool
	Focus() tea.Cmd
	Blur()
}
