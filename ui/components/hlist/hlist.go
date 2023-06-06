package hlist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	common       common.Common
	Style        Style
	Children     []common.Component
	offset       int
	active       int
	visibleItems int
}

type Style struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

func New(c common.Common, children ...common.Component) *Model {
	m := &Model{
		common: c,
		Style: Style{
			Normal: lipgloss.NewStyle(),
			Active: lipgloss.NewStyle(),
		},
		Children:     children,
		offset:       0,
		active:       0,
		visibleItems: c.Width, // set to highest possible ie 1 width items, set in Update()
	}

	if len(children) > 0 {
		switch current := m.Children[m.active].(type) {
		case common.FocusableComponent:
			current.Focus()
		}
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	// TODO should vlist and hlist have expanding children?
	// for _, child := range m.Children {
	// 	child.SetSize(height, child.Width())
	// }
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
		case key.Matches(msg, m.common.Global.KeyMap.Right):
			m.active = util.Min(m.active+1, len(m.Children)-1)

			if m.active == m.offset+m.visibleItems {
				m.offset = m.active
			}
		case key.Matches(msg, m.common.Global.KeyMap.Left):
			m.active = util.Max(m.active-1, 0)

			if m.active == m.offset-1 {
				m.offset = m.active
			}
		}

		if prevActive != m.active {
			switch prev := m.Children[prevActive].(type) {
			case common.FocusableComponent:
				prev.Blur()
			}

			switch current := m.Children[m.active].(type) {
			case common.FocusableComponent:
				current.Focus()
			}

		}
	}

	for _, child := range m.Children {
		_, cmd := child.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	width := 0
	m.visibleItems = 0

	var sections []string

	for i := m.offset; i < len(m.Children); i++ {
		section := m.Children[i].View()

		if i == m.active {
			section = m.Style.Active.Render(section)
		} else {
			section = m.Style.Normal.Render(section)
		}

		sectionWidth := lipgloss.Width(section)

		if width+sectionWidth > m.common.Width {
			break
		}

		width += sectionWidth
		m.visibleItems++

		sections = append(sections, section)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, sections...)
}
