package keymap

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit    key.Binding
	NextTab key.Binding
	PrevTab key.Binding
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Select  key.Binding
	Back    key.Binding
}

func NewKeyMap() *KeyMap {
	km := KeyMap{
		Quit:    key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit")),
		NextTab: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next tab")),
		PrevTab: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev tab")),
		Up:      key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		Down:    key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		Right:   key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "right")),
		Left:    key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "left")),
		Select:  key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select")),
		Back:    key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "back")),
	}

	return &km
}
