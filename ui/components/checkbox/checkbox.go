package checkbox

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

var (
	focusColor = lipgloss.Color("#F25D94")
	tabBorder  = lipgloss.RoundedBorder()

	borderStyle      = lipgloss.NewStyle().Border(tabBorder, true).BorderTop(false)
	focusBorderStyle = lipgloss.NewStyle().Border(tabBorder, true).BorderForeground(focusColor).Foreground(focusColor).BorderTop(false)

	topStyle     = lipgloss.NewStyle().Foreground(focusColor)
	checkedStyle = lipgloss.NewStyle().Foreground(focusColor)
)

type onChange func(value bool) tea.Cmd

type Model struct {
	props    common.Props
	focused  bool
	OnChange onChange
	Checked  bool
	Label    string
}

func New(p common.Props) *Model {
	return &Model{props: p, focused: false, Checked: false, OnChange: func(value bool) tea.Cmd { return nil }}
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
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Select):
			msg.Handled = true
			m.Checked = !m.Checked
			return m, m.OnChange(m.Checked)
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

	top := "╭───╮"
	if len(m.Label) == 4 {
		top = m.Label + "╮" //top[len(m.Label):] doesn't work b/c bytes
	}

	var bottom string
	if m.focused {
		top = topStyle.Render(top)
		bottom = focusBorderStyle.Render(pixel)
	} else {
		bottom = borderStyle.Render(pixel)
	}

	return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
}
