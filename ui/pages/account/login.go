package account

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

type LoginModel struct {
	common     common.Common
	httpClient *retryablehttp.Client
	inputs     []common.FocusableComponent
	focusIndex int
}

func NewLogin(c common.Common, httpClient *retryablehttp.Client) *LoginModel {
	m := &LoginModel{
		c,
		httpClient,
		make([]common.FocusableComponent, 2),
		0,
	}

	for i := range m.inputs {
		input := textinput.New()
		input.CursorStyle = cursorStyle
		input.CharLimit = 80
		switch i {
		case 0:
			input.Placeholder = "Email"
			input.Focus()
		case 1:
			input.Placeholder = "Password"
			input.EchoMode = textinput.EchoPassword
			// input.EchoCharacter = '*'
		}
		m.inputs[i] = &input
	}

	return m
}

func (m *LoginModel) SetSize(width, height int) {
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func blurFocusIndex(m *LoginModel) {
	// button focused
	if m.focusIndex == 2 {
		return
	}
	input := m.inputs[m.focusIndex].(*textinput.Model)
	input.Blur()
	input.PromptStyle = noStyle
	input.TextStyle = noStyle
}

func focusFocusIndex(m *LoginModel) {
	// button focused
	if m.focusIndex == 2 {
		return
	}
	input := m.inputs[m.focusIndex].(*textinput.Model)
	input.Focus()
	input.PromptStyle = focusedStyle
	input.TextStyle = focusedStyle
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case tea.KeyMsg:
		switch {
		// email, password, submit
		case key.Matches(msg, m.common.KeyMap.Down):
			blurFocusIndex(m)
			m.focusIndex = (m.focusIndex + 1) % 3
			focusFocusIndex(m)
		case key.Matches(msg, m.common.KeyMap.Up):
			blurFocusIndex(m)
			m.focusIndex = (m.focusIndex + 3 - 1) % 3
			focusFocusIndex(m)
		}
	}

	var cmd tea.Cmd

	for i := range m.inputs {
		switch input := m.inputs[i].(type) {
		case *textinput.Model:
			var f textinput.Model
			f, cmd = input.Update(msg)
			m.inputs[i] = &f
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *LoginModel) View() string {
	var sections []string

	for i := range m.inputs {
		sections = append(sections, m.inputs[i].View())
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
