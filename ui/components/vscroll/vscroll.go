package vscroll

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	props    common.Props
	children []tea.Model
	offset   int
}

func New(p common.Props, children []tea.Model) *Model {
	return &Model{
		props:    p,
		children: children,
		offset:   0,
	}
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.props.Global.KeyMap.Down):
			m.offset++
			// TODO how do upper bound for dynamic height

		case key.Matches(msg, m.props.Global.KeyMap.Up):
			m.offset = util.Max(m.offset-1, 0)
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

	started := false
	height := 0

	for _, child := range m.children {
		heightLeft := m.props.Height + m.offset - height

		if heightLeft <= 0 {
			break
		}

		section := child.View()
		sectionHeight := lipgloss.Height(section)

		height += sectionHeight

		if !started {
			if height >= m.offset {
				started = true
			}
			if height > m.offset {
				something := height - m.offset
				subSections := strings.Split(section, "\n")
				// TODO if taller than viewport height
				visiblePart := subSections[(sectionHeight - something):]
				sb.WriteString(strings.Join(visiblePart, "\n"))
			}
			continue
		}

		if heightLeft < sectionHeight {
			subSections := strings.SplitN(section, "\n", heightLeft+1)
			visiblePart := subSections[:len(subSections)-1]

			sb.WriteString(strings.Join(visiblePart, "\n"))
			continue
		}

		sb.WriteString(section)

		if height < m.props.Height {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
