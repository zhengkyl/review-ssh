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
	Style        Style
	children     []common.Component
	offset       int
	active       int
	visibleItems int
}

type Style struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

func New(c common.Common, children []common.Component) *Model {
	m := &Model{
		common: c,
		Style: Style{
			Normal: lipgloss.NewStyle(),
			Active: lipgloss.NewStyle(),
		},
		children:     children,
		offset:       0,
		active:       0,
		visibleItems: c.Height, // set to highest possible ie 1 height items, set in Update()
	}

	if len(children) > 0 {
		switch current := m.children[m.active].(type) {
		case common.FocusableComponent:
			current.Focus()
		}
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	for _, child := range m.children {
		child.SetSize(width, child.Height())
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		prevActive := m.active
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

		if prevActive != m.active {
			switch prev := m.children[prevActive].(type) {
			case common.FocusableComponent:
				prev.Blur()
			}

			switch current := m.children[m.active].(type) {
			case common.FocusableComponent:
				current.Focus()
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

	for i := m.offset; i < len(m.children); i++ {
		section := m.children[i].View()

		if i == m.active {
			section = m.Style.Active.Render(section)
		} else {
			section = m.Style.Normal.Render(section)
		}

		sectionHeight := lipgloss.Height(section)

		if height+sectionHeight > m.common.Height {
			break
		}

		height += sectionHeight
		m.visibleItems++

		if i > m.offset {
			sb.WriteString("\n")
		}

		sb.WriteString(section)
	}

	return sb.String()
}
