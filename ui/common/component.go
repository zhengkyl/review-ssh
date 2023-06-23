package common

// "github.com/charmbracelet/bubbles/help"

type Sizable interface {
	SetSize(width, height int)
	Width() int
	Height() int
	// GetMargins() int
}

type Component interface {
	Model
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
	Focus()
	Blur()
}
