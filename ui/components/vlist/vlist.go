package vlist

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	common       common.Common
	children     []tea.Model
	offset       int
	active       int
	visibleItems int
}

func New(c common.Common, children []tea.Model) *Model {
	return &Model{
		common:       c,
		children:     children,
		offset:       0,
		active:       0,
		visibleItems: 1, // ??? TODO
	}
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.common.Global.KeyMap.Down):
			m.active = util.Min(m.active+1, len(m.children)-1)

			if m.active == m.offset+m.visibleItems {
				m.offset = m.active
			}
		case key.Matches(msg, m.common.Global.KeyMap.Up):
			m.active = util.Max(m.active-1, 0)

			if m.active == m.offset-1 {
				m.offset = m.active
			}
		}
	}
	for _, child := range m.children {
		_, cmd := child.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}
	height := 0

	m.visibleItems = 0
	// sb.WriteString(fmt.Sprintf("OFFSET %v : ACTIVE %v \n", m.offset, m.active))

	for i := m.offset; i < len(m.children); i++ {
		section := m.children[i].View()
		sectionHeight := lipgloss.Height(section)

		if height+sectionHeight > m.common.Height {
			break
		}

		height += sectionHeight
		m.visibleItems++

		sb.WriteString(section)
		sb.WriteString("\n")
	}

	return sb.String()
}
