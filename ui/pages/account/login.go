package account

import (
	"bytes"
	"encoding/json"

	"github.com/charmbracelet/bubbles/key"
	bubbles_textinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
	"github.com/zhengkyl/review-ssh/ui/components/textinput"
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
	state      *common.Shared
	inputs     []common.FocusableComponent
	focusIndex int
}

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLogin(c common.Common, s *common.Shared) *LoginModel {
	m := &LoginModel{
		c,
		s,
		make([]common.FocusableComponent, 3),
		0,
	}

	for i := 0; i < SUBMIT_INDEX; i++ {
		input := textinput.New(c)
		input.Inner.CursorStyle = cursorStyle
		input.Inner.CharLimit = 80
		switch i {
		case 0:
			input.Inner.Placeholder = "Email"
			input.Focus()
		case 1:
			input.Inner.Placeholder = "Password"
			input.Inner.EchoMode = bubbles_textinput.EchoPassword
			// input.EchoCharacter = '*'
		}
		m.inputs[i] = input
	}

	m.inputs[2] = button.New(c, "Submit", func() tea.Msg { return nil })

	return m
}

func (m *LoginModel) SetSize(width, height int) {
}

func (m *LoginModel) Init() tea.Cmd {
	return nil
}

func blurFocusIndex(m *LoginModel) {
	// button focused
	m.inputs[m.focusIndex].Blur()
	if m.focusIndex == SUBMIT_INDEX {
		return
	}
	input := m.inputs[m.focusIndex].(*textinput.Model)
	input.Inner.PromptStyle = noStyle
	input.Inner.TextStyle = noStyle
}

func focusFocusIndex(m *LoginModel) {
	// button focused
	m.inputs[m.focusIndex].Focus()
	if m.focusIndex == SUBMIT_INDEX {
		return
	}
	input := m.inputs[m.focusIndex].(*textinput.Model)
	input.Inner.PromptStyle = focusedStyle
	input.Inner.TextStyle = focusedStyle
}

func changeFocusIndex(m *LoginModel, newIndex int) {
	blurFocusIndex(m)
	m.focusIndex = newIndex
	focusFocusIndex(m)
}

func postAuth(client *retryablehttp.Client, loginData loginData) tea.Cmd {
	return func() tea.Msg {

		bsLoginData, err := json.Marshal(loginData)

		if err != nil {
			return common.AuthState{
				Authed: false,
			}
		}

		resp, err := client.Post("https://review-api.fly.dev/auth", "application/json", bytes.NewBuffer(bsLoginData))

		if err != nil {
			return common.AuthState{
				Authed: false,
			}
		}

		if resp.StatusCode != 204 {
			return common.AuthState{
				Authed: false,
			}
		}

		cookie := resp.Header.Get("Set-Cookie")

		return common.AuthState{
			Authed: true,
			Cookie: cookie,
		}
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

				return m, postAuth(&m.state.HttpClient, loginData{
					m.inputs[0].(*textinput.Model).Inner.Value(),
					m.inputs[1].(*textinput.Model).Inner.Value(),
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
		_, cmd = m.inputs[i].Update(msg)

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *LoginModel) View() string {
	var sections []string

	for i := range m.inputs {
		sections = append(sections, m.inputs[i].View())
	}

	sections = append(sections, m.state.AuthState.Cookie)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
