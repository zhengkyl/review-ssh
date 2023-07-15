package vlist

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type overflow int

const (
	Scroll overflow = iota
	Paginate
)

type Model struct {
	props      common.Props
	Style      Style
	items      []common.Focusable
	offset     int
	active     int
	perPage    int
	ItemHeight int
	ItemGap    int
	Overflow   overflow
}

type Style struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

func New(p common.Props, itemHeight int, items ...common.Focusable) *Model {
	m := &Model{
		props: p,
		Style: Style{
			Normal: lipgloss.NewStyle(),
			Active: lipgloss.NewStyle(),
		},
		items:      items,
		offset:     0,
		active:     0,
		ItemHeight: itemHeight,
		ItemGap:    1,
		Overflow:   Scroll,
	}

	m.SetSize(p.Width, p.Height)

	if len(items) > 0 {
		current := m.items[m.active]
		current.Focus()
	}

	return m
}

func (m *Model) Offset() int {
	return m.offset
}

func (m *Model) PerPage() int {
	return m.perPage
}

func (m *Model) Length() int {
	return len(m.items)
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	for _, child := range m.items {
		sizable, ok := child.(common.Sizable)
		if ok {
			sizable.SetSize(width, m.ItemHeight)
		}
	}

	m.perPage = (m.props.Height + m.ItemGap) / (m.ItemHeight + m.ItemGap)

	switch m.Overflow {
	case Scroll:
		// Try to keep active item same pos from top when resizing
		maxIndex := util.Max(m.perPage-1, 0)
		newIndex := util.Min(m.active-m.offset, maxIndex)
		m.offset = m.active - newIndex
	case Paginate:
		m.offset = m.active / util.Max(m.perPage, 1)
	}
}

func (m *Model) SetItems(items []common.Focusable) {
	m.items = items
	m.active = 0
	if len(items) > 0 {
		current := m.items[m.active]
		current.Focus()
	}
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case *common.KeyEvent:
		_, cmd := m.items[m.active].Update(msg)

		if msg.Handled {
			return m, cmd
		}

		prevActive := m.active
		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Down):
			msg.Handled = true
			m.active = util.Min(m.active+1, len(m.items)-1)

			if m.active == m.offset+m.perPage {
				switch m.Overflow {
				case Scroll:
					m.offset++
				case Paginate:
					m.offset += m.perPage
				}

			}
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Up):
			msg.Handled = true
			m.active = util.Max(m.active-1, 0)

			if m.active == m.offset-1 {
				switch m.Overflow {
				case Scroll:
					m.offset = m.active
				case Paginate:
					m.offset -= m.perPage
				}
			}
		}

		if prevActive != m.active {
			prev := m.items[prevActive]
			prev.Blur()

			current := m.items[m.active]
			current.Focus()

		}
	}

	for _, child := range m.items {
		_, cmd := child.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}

	for i := m.offset; i < m.offset+m.perPage && i < len(m.items); i++ {
		section := m.items[i].View()

		if i == m.active {
			section = m.Style.Active.Render(section)
		} else {
			section = m.Style.Normal.Render(section)
		}

		if i > m.offset {
			sb.WriteString(strings.Repeat("\n", m.ItemGap+1))
		}

		sb.WriteString(section)
	}

	return sb.String()
}
