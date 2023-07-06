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
	props        common.Props
	Style        Style
	Children     []common.Focusable
	offset       int
	Active       int
	visibleItems int
}

type Style struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

func New(p common.Props, children ...common.Focusable) *Model {
	m := &Model{
		props: p,
		Style: Style{
			Normal: lipgloss.NewStyle(),
			Active: lipgloss.NewStyle(),
		},
		Children:     children,
		offset:       0,
		Active:       0,
		visibleItems: p.Height, // set to highest possible ie 1 height items, set in Update()
	}

	if len(children) > 0 {
		current := m.Children[m.Active]
		current.Focus()
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	for _, child := range m.Children {
		sizable, ok := child.(common.Sizable)
		if ok {
			sizable.SetSize(width, sizable.Height())
		}
	}

	// Try to keep active item same pos from top when resizing
	maxIndex := util.Max(m.visibleItems-1, 0)
	newIndex := util.Min(m.Active-m.offset, maxIndex)
	m.offset = m.Active - newIndex
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case *common.KeyEvent:
		_, cmd := m.Children[m.Active].Update(msg)
		if msg.Handled {
			return m, cmd
		}

		prevActive := m.Active
		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Down):
			m.Active = util.Min(m.Active+1, len(m.Children)-1)

			if m.Active == m.offset+m.visibleItems {
				m.offset++
			}
			msg.Handled = true
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Up):
			m.Active = util.Max(m.Active-1, 0)

			if m.Active == m.offset-1 {
				m.offset = m.Active
			}
			msg.Handled = true
		}

		if prevActive != m.Active {
			prev := m.Children[prevActive]
			prev.Blur()

			current := m.Children[m.Active]
			current.Focus()

		}
	}

	for _, child := range m.Children {
		_, cmd := child.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}
	height := 0

	m.visibleItems = 0

	for i := m.offset; i < len(m.Children); i++ {
		section := m.Children[i].View()

		if i == m.Active {
			section = m.Style.Active.Render(section)
		} else {
			section = m.Style.Normal.Render(section)
		}

		sectionHeight := lipgloss.Height(section)

		if height+sectionHeight > m.props.Height {
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
