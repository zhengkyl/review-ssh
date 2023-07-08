package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ShowFilm int

type KeyEvent struct {
	KeyMsg  tea.KeyMsg
	Handled bool
}
