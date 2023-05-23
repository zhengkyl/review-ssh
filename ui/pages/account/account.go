package account

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

const SUBMIT_INDEX = 2

type Model struct {
	common     common.Common
	inputs     []common.FocusableComponent
	focusIndex int
}

func New(c common.Common) *Model {
	m := &Model{
		c,
		make([]common.FocusableComponent, 3),
		0,
	}

	inputCommon := common.Common{
		Width:  c.Width - 0, // TODO padding
		Height: 3,           // TODO does nothing
		Global: c.Global,
	}

	for i := 0; i < SUBMIT_INDEX; i++ {
		input := textfield.New(inputCommon)
		input.CursorStyle(cursorStyle)
		input.CharLimit(80)

		switch i {
		case 0:
			input.Placeholder("Email")
			input.Focus()
		case 1:
			input.Placeholder("Password")
			input.EchoMode(textinput.EchoPassword)
			// input.EchoCharacter = '*'
		}
		m.inputs[i] = input
	}

	m.inputs[2] = button.New(c, "Submit", func() tea.Msg { return nil })

	return m
}

func (m *Model) SetSize(width, height int) {
	m.inputs[0].SetSize(width, 3)
	m.inputs[1].SetSize(width, 3)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func blurFocusIndex(m *Model) {
	// button focused
	m.inputs[m.focusIndex].Blur()
	if m.focusIndex == SUBMIT_INDEX {
		return
	}
	input := m.inputs[m.focusIndex].(*textfield.Model)
	input.PromptStyle(noStyle)
	input.TextStyle(noStyle)
}

func focusFocusIndex(m *Model) {
	// button focused
	m.inputs[m.focusIndex].Focus()
	if m.focusIndex == SUBMIT_INDEX {
		return
	}
	input := m.inputs[m.focusIndex].(*textfield.Model)
	input.PromptStyle(focusedStyle)
	input.TextStyle(focusedStyle)
}

func changeFocusIndex(m *Model, newIndex int) {
	blurFocusIndex(m)
	m.focusIndex = newIndex
	focusFocusIndex(m)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		for _, input := range m.inputs {
			input.SetSize(msg.Width, 3)
		}

	case tea.KeyMsg:
		switch {
		// email, password, submit
		case key.Matches(msg, m.common.Global.KeyMap.Select):
			if m.focusIndex == SUBMIT_INDEX {

				return m, postAuth(&m.common.Global.HttpClient, loginData{
					m.inputs[0].(*textfield.Model).Value(),
					m.inputs[1].(*textfield.Model).Value(),
				})
			}
			changeFocusIndex(m, (m.focusIndex+1)%3)
		case key.Matches(msg, m.common.Global.KeyMap.NextInput):
			changeFocusIndex(m, (m.focusIndex+1)%3)
		case key.Matches(msg, m.common.Global.KeyMap.PrevInput):
			changeFocusIndex(m, (m.focusIndex+3-1)%3)
		}
	}

	var cmd tea.Cmd
	for i := range m.inputs {
		_, cmd = m.inputs[i].Update(msg)

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// return "-0\n1\n2\n3\n4\n5\n6\n7\n8\n-9"

	// if m.global.AuthState.Authed {
	// 	return m.global.AuthState.User.Name
	// }

	var sections []string

	sections = append(sections, m.common.Global.AuthState.Cookie)

	for i := range m.inputs {
		sections = append(sections, m.inputs[i].View())
	}

	// sections = append(sections, m.global.AuthState.Cookie)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
