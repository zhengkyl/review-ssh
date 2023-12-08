package search

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/components/vlist"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	viewStyle = lipgloss.NewStyle().MarginTop(1)
)

type Model struct {
	props       common.Props
	list        *vlist.Model
	searchField *textfield.Model
	focused     bool
	Query       string
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

	vf := viewStyle.GetVerticalFrameSize()

	m.list.SetSize(width, height-vf-1)
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	_, cmd := m.list.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	sb := strings.Builder{}
	if m.list.Length() == 0 {
		sb.WriteString("No results.")
	} else {
		sb.WriteString(m.list.View())
	}

	viewH := m.props.Height - 1
	sb.WriteString(strings.Repeat("\n", viewH-lipgloss.Height(sb.String())))

	start := m.list.Offset() + 1
	last := m.list.Offset() + m.list.PerPage()
	paginator := fmt.Sprintf("%d-%d of %d", start, util.Max(last, m.list.Length()), m.list.Length())

	// filmitems have an internal horizontal framesize of 2
	sb.WriteString(strings.Repeat(" ", m.props.Width-len(paginator)-2))
	sb.WriteString(paginator)

	return viewStyle.Render(sb.String())
}
