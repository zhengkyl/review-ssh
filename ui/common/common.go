package common

import (
	"github.com/zhengkyl/review-ssh/ui/keymap"
	"github.com/zhengkyl/review-ssh/ui/styles"
)

type Common struct {
	Width  int
	Height int
	Styles *styles.Styles
	KeyMap *keymap.KeyMap
}

func (c *Common) SetSize(width, height int) {
	c.Width = width
	c.Height = height
}
