package button

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Model struct {
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

func New(common common.Common, text string, callback tea.Cmd) *Model {
	return &Model{
		common,
		text,
		callback,
		false,
	}
}

func (m *Model) Focused() bool {
	return m.focus
}

func (m *Model) Focus() tea.Cmd {
	m.focus = true
	return nil
}

func (m *Model) Blur() {
	m.focus = false
}

func (m *Model) SetSize(h, w int) {

}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.focus {

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.common.Global.KeyMap.Select):
				return m, m.callback
			}
		}

	}

	return m, nil
}

func (m *Model) View() string {
	if m.focus {
		return activeButtonStyle.Render(m.text)
	}

	return buttonStyle.Render(m.text)
}
