package common

import "github.com/zhengkyl/review-ssh/ui/common/enums"

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
