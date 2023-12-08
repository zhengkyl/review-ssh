package keymap

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit   key.Binding
	Help   key.Binding
	Search key.Binding
	NextX  key.Binding
	PrevX  key.Binding
	NextY  key.Binding
	PrevY  key.Binding
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Back   key.Binding
	Move   key.Binding
}

func DefaultKeyMap() *KeyMap {
	km := KeyMap{
		Quit:   key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:   key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Search: key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "search")),
		NextX:  key.NewBinding(key.WithKeys("tab", "right", "l"), key.WithHelp("tab", "next tab")),
		PrevX:  key.NewBinding(key.WithKeys("shift+tab", "left", "h"), key.WithHelp("shift+tab", "prev tab")),
		// TODO temp fix enter is NextY b/c input lists
		NextY:  key.NewBinding(key.WithKeys("tab", "down", "j", "enter"), key.WithHelp("tab", "next tab")),
		PrevY:  key.NewBinding(key.WithKeys("shift+tab", "up", "k"), key.WithHelp("shift+tab", "prev tab")),
		Up:     key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
		Down:   key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
		Right:  key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "right")),
		Left:   key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "left")),
		Select: key.NewBinding(key.WithKeys("enter", " "), key.WithHelp("[enter]", "select")),
		Back:   key.NewBinding(key.WithKeys("esc"), key.WithHelp("[esc]", "back")),
		//
		Move: key.NewBinding(key.WithKeys("right", "down", "up", "left"), key.WithHelp("←↓↑→", "move")),
	}

	return &km
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Quit, k.Move, k.Select, k.Back}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		// {k.Up, k.Down, k.Left, k.Right}, // first column
		// {k.Help, k.Quit},                // second column
	}
}
