package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Model struct {
	Inner textinput.Model
}

func New(common common.Common) *Model {
	inner := textinput.New()

	return &Model{inner}
}

func (m *Model) Focus() tea.Cmd {
	return m.Inner.Focus()
}

func (m *Model) Blur() {
	m.Inner.Blur()
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Inner, cmd = m.Inner.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	// TODO padding etc
	return m.Inner.View()
}
