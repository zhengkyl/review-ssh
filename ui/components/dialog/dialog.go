package dialog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/hlist"
)

type Model struct {
	common  common.Common
	text    string
	buttons *hlist.Model
}

func New(c common.Common) *Model {
	return &Model{
		common:  c,
		buttons: hlist.New(c),
	}
}

func (m *Model) SetSize(h, w int) {

}

func (m *Model) Height() int {
	return m.common.Height
}

func (m *Model) Width() int {
	return m.common.Width
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) View() string {
	return ""
}
