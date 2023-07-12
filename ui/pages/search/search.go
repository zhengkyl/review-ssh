package search

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/components/vlist"
)

type Model struct {
	props       common.Props
	list        *vlist.Model
	searchField *textfield.Model
	focused     bool
}

func New(p common.Props, searchField *textfield.Model) *Model {
	m := &Model{
		props:       p,
		list:        vlist.New(p, 6),
		searchField: searchField,
		focused:     false,
	}

	m.list.Overflow = vlist.Paginate

	m.SetSize(p.Width, p.Height)

	return m
}

func (m *Model) SetItems(items []common.Focusable) {
	m.list.SetItems(items)
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	m.list.SetSize(width, height)
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {

	var cmds []tea.Cmd

	// switch msg := msg.(type) {
	// case *common.KeyEvent:
	// 	switch {
	// 	}
	// }

	var cmd tea.Cmd

	_, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}
	// ss := lipgloss.NewStyle().Width(m.props.Width - wm)
	// view = ss.Render(m.input.View())
	sb.WriteString(m.list.View())

	return sb.String()
}
