package common

import "github.com/zhengkyl/review-ssh/ui/styles"

type Common struct {
	Width  int
	Height int
	Styles *styles.Styles
}

func (c *Common) SetSize(width, height int) {
	c.Width = width
	c.Height = height
}
