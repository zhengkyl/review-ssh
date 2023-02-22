package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type ButtonModel struct {
	common   common.Common
	text     string
	callback tea.Cmd
	focus    bool
}

var (
	buttonStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Background(lipgloss.Color("#888B7E")).Padding(0, 1)
	activeButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).Padding(0, 1)
)

func New(common common.Common, text string, callback tea.Cmd) *ButtonModel {
	return &ButtonModel{
		common,
		text,
		callback,
		false,
	}
}

func (m *ButtonModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	return tea.Batch(cmds...)
}

func (m *ButtonModel) Focus() tea.Cmd {
	m.focus = true
	return nil
}

func (m *ButtonModel) Blur() {
	m.focus = false
}

func (m *ButtonModel) Update(msg tea.Msg) (*ButtonModel, tea.Cmd) {
	return m, nil

	// var cmds []tea.Cmd

	// // m.buttons[m.activeButton].text = "no"

	// return m, tea.Batch(cmds...)
}

func (m *ButtonModel) View() string {
	if m.focus {
		return activeButtonStyle.Render(m.text)
	}

	return buttonStyle.Render(m.text)
}
