package common

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
)

type ShowPage struct {
	Category enums.Category
	Tmdb_id  int
	Season   int
}

type KeyEvent struct {
	KeyMsg  tea.KeyMsg
	Handled bool
}
