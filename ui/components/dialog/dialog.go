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
	common  common.Common
	text    string
	buttons []button.Model
	active  int
	focused bool
}

func New(c common.Common, text string) *Model {
	m := &Model{
		common: c,
		text:   text,
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

func (m *Model) Focus() tea.Cmd {
	if m.active != 0 {
		m.buttons[m.active].Blur()
		m.active = 0
		m.buttons[m.active].Focus()
	}

	m.focused = true
	return nil
}

func (m *Model) Blur() {
	m.focused = false
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
		case key.Matches(msg, m.common.Global.KeyMap.NextInput):
			m.active = util.Min(m.active+1, len(m.buttons)-1)
		case key.Matches(msg, m.common.Global.KeyMap.PrevInput):
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
	sb.WriteString("\n\n")
	for _, button := range m.buttons {
		sb.WriteString(button.View())
		sb.WriteString(" ")
	}
	return dialogStyle.Render(sb.String())
}