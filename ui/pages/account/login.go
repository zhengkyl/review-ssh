package account

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle  = focusedStyle.Copy()
	noStyle      = lipgloss.NewStyle()
)

const SUBMIT_INDEX = 2

type LoginModel struct {
	common     common.Common
	httpClient *retryablehttp.Client
	inputs     []common.FocusableComponent
	focusIndex int
}

type loginData struct {
	email    string
	password string
}

func NewLogin(c common.Common, httpClient *retryablehttp.Client) *LoginModel {
	m := &LoginModel{
		c,
		httpClient,
		make([]common.FocusableComponent, 3),
		0,
	}

	for i := 0; i < SUBMIT_INDEX; i++ {
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

	m.inputs[2] = button.New(c, "Submit", func() tea.Msg { return nil })

	return m
}

func (m *LoginModel) SetSize(width, height int) {
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func blurFocusIndex(m *LoginModel) {
	// button focused
	m.inputs[m.focusIndex].Blur()
	if m.focusIndex == SUBMIT_INDEX {
		return
	}
	input := m.inputs[m.focusIndex].(*textinput.Model)
	input.PromptStyle = noStyle
	input.TextStyle = noStyle
}

func focusFocusIndex(m *LoginModel) {
	// button focused
	m.inputs[m.focusIndex].Focus()
	if m.focusIndex == SUBMIT_INDEX {
		return
	}
	input := m.inputs[m.focusIndex].(*textinput.Model)
	input.PromptStyle = focusedStyle
	input.TextStyle = focusedStyle
}

func changeFocusIndex(m *LoginModel, newIndex int) {
	blurFocusIndex(m)
	m.focusIndex = newIndex
	focusFocusIndex(m)
}

func postAuth(client *retryablehttp.Client, loginData loginData) tea.Cmd {
	return func() tea.Msg {

		resp, err := client.Post("https://review-api.fly.dev/auth", "application/json", loginData)

		if err != nil {
			return nil
		}

		if resp.StatusCode != 204 {
			return nil
		}

		cookie := resp.Header.Get("Set-Cookie")

		return cookie
	}
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case tea.KeyMsg:
		switch {
		// email, password, submit
		case key.Matches(msg, m.common.KeyMap.Select):
			if m.focusIndex == SUBMIT_INDEX {

				return m, postAuth(m.httpClient, loginData{
					m.inputs[0].(*textinput.Model).Value(),
					m.inputs[1].(*textinput.Model).Value(),
				})
			}
			changeFocusIndex(m, (m.focusIndex+1)%3)
		case key.Matches(msg, m.common.KeyMap.NextInput):
			changeFocusIndex(m, (m.focusIndex+1)%3)
		case key.Matches(msg, m.common.KeyMap.PrevInput):
			changeFocusIndex(m, (m.focusIndex+3-1)%3)
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
