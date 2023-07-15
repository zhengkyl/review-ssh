package search

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/components/vlist"
)

var (
	viewStyle = lipgloss.NewStyle().Margin(1)
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
	hf := viewStyle.GetHorizontalFrameSize()
	vf := viewStyle.GetVerticalFrameSize()

	m.props.Width = width - hf
	m.props.Height = height - vf

	m.list.SetSize(width, height-2)
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

	sb.WriteString("\n\n")
	perPage := m.list.PerPage()
	page := m.list.Offset() / perPage
	pages := (m.list.Length() + perPage - 1) / perPage
	paginator := fmt.Sprintf("%d/%d", page, pages)
	sb.WriteString(strings.Repeat(" ", (m.props.Width/2)-len(paginator)))
	sb.WriteString(paginator)

	return viewStyle.Render(sb.String())
}
