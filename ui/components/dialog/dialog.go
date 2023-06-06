package dialog

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	common  common.Common
	text    string
	buttons []button.Model
	active  int
}

func New(c common.Common, text string, buttons ...button.Model) *Model {
	m := &Model{
		common:  c,
		text:    text,
		buttons: buttons,
		active:  0,
	}

	m.buttons[0].Focus()

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
}

func (m *Model) Height() int {
	return m.common.Height
}

func (m *Model) Width() int {
	return m.common.Width
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
			m.active = util.Min(m.active+1, len(m.buttons)-1)
		case key.Matches(msg, m.common.Global.KeyMap.Left):
			m.active = util.Max(m.active-1, 0)
		}

		if prevActive != m.active {
			m.buttons[prevActive].Blur()
			m.buttons[m.active].Focus()
		}
	}

	for _, child := range m.buttons {
		_, cmd := child.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)

}

func (m *Model) View() string {
	sb := strings.Builder{}

	sb.WriteString(m.text)
	sb.WriteString("\n")
	for _, button := range m.buttons {
		sb.WriteString(button.View())
	}
	sb.WriteString("\n")
	return sb.String()
}
