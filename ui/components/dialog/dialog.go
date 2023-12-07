package dialog

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	dialogStyle = lipgloss.NewStyle().Padding(1, 3).Border(lipgloss.RoundedBorder(), true)
)

type Model struct {
	props   common.Props
	text    string
	buttons []button.Model
	active  int
	focused bool
}

func New(p common.Props, text string) *Model {
	m := &Model{
		props: p,
		text:  text,
	}

	return m
}

func (m *Model) Buttons(buttons ...button.Model) {
	m.buttons = buttons

	m.active = 0
	m.buttons[0].Focus()
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	if m.active != 0 {
		m.buttons[m.active].Blur()
		m.active = 0
		m.buttons[m.active].Focus()
	}

	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case *common.KeyEvent:
		prevActive := m.active

		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Back):
			msg.Handled = true
			m.Blur()
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.NextX):
			msg.Handled = true
			m.active = util.Mod(m.active+1, len(m.buttons))
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.PrevX):
			msg.Handled = true
			m.active = util.Mod(m.active-1, len(m.buttons))
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
	sb.WriteString("\n\n")
	for _, button := range m.buttons {
		sb.WriteString(button.View())
		sb.WriteString(" ")
	}
	return dialogStyle.Render(sb.String())
}
