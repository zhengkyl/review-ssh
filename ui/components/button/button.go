package button

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Model struct {
	props    common.Props
	Style    Style
	text     string
	callback tea.Cmd
	focused  bool
}

type Style struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

func New(p common.Props, text string, callback tea.Cmd) *Model {
	return &Model{
		props: p,
		Style: Style{
			Normal: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Background(lipgloss.Color("#888B7E")).Padding(0, 1),
			Active: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFF7DB")).Background(lipgloss.Color("#F25D94")).Padding(0, 1),
		},
		text:     text,
		callback: callback,
		focused:  false,
	}
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

func (m *Model) SetSize(h, w int) {

}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	if m.focused {

		switch msg := msg.(type) {
		case *common.KeyEvent:
			switch {
			case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Select):
				msg.Handled = true
				return m, m.callback
			}
		}

	}

	return m, nil
}

func (m *Model) View() string {
	if m.focused {
		return m.Style.Active.Render(m.text)
	}

	return m.Style.Normal.Render(m.text)
}
