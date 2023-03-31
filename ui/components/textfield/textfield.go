package textfield

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

var (
	tabBorder  = lipgloss.RoundedBorder()
	inputStyle = lipgloss.NewStyle().Border(tabBorder, true)
)

type Model struct {
	common common.Common
	Inner  textinput.Model
}

func New(common common.Common) *Model {
	inner := textinput.New()
	inner.Width = common.Width

	return &Model{common, inner}
}

func (m *Model) Focus() tea.Cmd {
	return m.Inner.Focus()
}

func (m *Model) Blur() {
	m.Inner.Blur()
}

func (m *Model) SetSize(w, h int) {
	// TODO figure out how to get this padding/margin

	m.common.Width = w - 8
	m.common.Height = h

	m.Inner.Width = w - 8

	// TODO what do if background color set
	if m.Inner.Placeholder != "" {
		m.Inner.Placeholder = m.Inner.Placeholder + strings.Repeat(" ", m.Inner.Width-len(m.Inner.Placeholder))
	}
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
	return inputStyle.Render(m.Inner.View())
}
