package checkbox

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

var (
	tabBorder        = lipgloss.RoundedBorder()
	borderStyle      = lipgloss.NewStyle().Border(tabBorder, true)
	focusBorderStyle = lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(lipgloss.Color("227"))
	checkedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("227"))
)

type Model struct {
	props    common.Props
	focused  bool
	Callback tea.Cmd
	Checked  bool
}

func New(p common.Props) *Model {
	return &Model{props: p, focused: false, Checked: false}
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *common.KeyEvent:
		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Back):
			m.Blur()
			msg.Handled = true
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Select):
			m.Checked = !m.Checked
			msg.Handled = true
		}
	}
	return m, nil
}

func (m *Model) View() string {
	var pixel string
	if m.Checked {
		pixel = checkedStyle.Render("▐█▌")
	} else {
		pixel = "   "
	}

	if m.focused {
		return focusBorderStyle.Render(pixel)
	} else {
		return borderStyle.Render(pixel)
	}
}
