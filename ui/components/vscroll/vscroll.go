package vertical

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Model struct {
	common   common.Common
	children []tea.Model
}

func New(c common.Common, children []tea.Model) *Model {
	return &Model{
		common:   c,
		children: children,
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
	for _, child := range m.children {
		_, cmd := child.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}
	height := 0

	for i, child := range m.children {
		if m.common.Height == height {
			break
		}

		section := child.View()
		sectionHeight := lipgloss.Height(section)

		heightLeft := m.common.Height - (height + sectionHeight)

		if heightLeft <= 0 {
			subSections := strings.SplitN(section, "\n", -heightLeft+1)
			visiblePart := subSections[:len(subSections)-1]
			sb.WriteString(strings.Join(visiblePart, "\n"))
			break
		}

		sb.WriteString(section)
		if i != len(m.children)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
