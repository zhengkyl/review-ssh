package common

import tea "github.com/charmbracelet/bubbletea"

// "github.com/charmbracelet/bubbles/help"
type Model interface {
	// Update is called when a message is received. Use it to inspect messages
	// and, in response, update the model and/or send a command.
	Update(tea.Msg) (Model, tea.Cmd)

	// View renders the program's UI, which is just a string. The view is
	// rendered after every Update.
	View() string
}

type Props struct {
	Width  int
	Height int
	Global Global
}
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
