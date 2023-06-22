package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
)

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

type ShowPage struct {
	Category enums.Category
	Tmdb_id  int
	Season   int
}

// func (c *Common) SetSize(width, height int) {
// 	c.Width = width
// 	c.Height = height
// }
