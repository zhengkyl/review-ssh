package common

import (
	// "github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

// no Init() b/c bubble tea components don't implement it
// no Update() b/c it makes sense to return a typed pointer, not tea.Model
type CustomModel interface {
	View() string
}

type Sizable interface {
	SetSize(width, height int)
}

type PageComponent interface {
	tea.Model
	Sizable
	// help.KeyMap
}

type Component interface {
	CustomModel
	Sizable
}

type FocusableComponent interface {
	CustomModel
	Focus() tea.Cmd
	Blur()
}
