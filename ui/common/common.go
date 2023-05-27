package common

import "github.com/charmbracelet/lipgloss"

type Common struct {
	Width  int
	Height int
	Style  lipgloss.Style
	Global Global
}

// func (c *Common) SetSize(width, height int) {
// 	c.Width = width
// 	c.Height = height
// }
