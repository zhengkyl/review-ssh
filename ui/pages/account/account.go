package account

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/components/vlist"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle  = focusedStyle.Copy()
)

const SUBMIT_INDEX = 2

type Model struct {
	common     common.Common
	inputs     *vlist.Model
	buttons    *vlist.Model
	focusIndex int
	stage      int
	debug      string
}

const (
	picker = 0
	signIn = 2
	signUp = 3
)

type signInMsg struct{}
type signUpMsg struct{}

func New(c common.Common) *Model {
	b := vlist.New(c,
		button.New(c, "     Sign in     ", func() tea.Msg { return signInMsg{} }),
		button.New(c, "     Sign up     ", func() tea.Msg { return signUpMsg{} }),
		button.New(c, "Continue as guest", func() tea.Msg { return common.GuestAuthState }),
	)

	b.Style.Active = lipgloss.NewStyle().MarginTop(1)
	b.Style.Normal = lipgloss.NewStyle().MarginTop(1)

	m := &Model{
		common:     c,
		inputs:     vlist.New(c),
		buttons:    b,
		focusIndex: 0,
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	m.inputs.SetSize(width, height)
	m.buttons.SetSize(width, height)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case signUpRes:
		m.debug = fmt.Sprint(msg)
	case signInMsg:
		m.stage = signIn
		m.inputs.Children = signInInputs(m.common)
	case signUpMsg:
		m.stage = signUp
		m.inputs.Children = signUpInputs(m.common)
	case tea.KeyMsg:
		switch {
		}
	}

	var cmd tea.Cmd

	if m.stage == picker {
		_, cmd = m.buttons.Update(msg)
	} else {
		_, cmd = m.inputs.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// return "-0\n1\n2\n3\n4\n5\n6\n7\n8\n-9"
	// return fmt.Sprintf("ITEM: %p\n1\n2\n3\n4\n5\n6\n7", m)

	// if m.global.AuthState.Authed {
	// 	return m.global.AuthState.User.Name
	// }

	sb := strings.Builder{}

	sb.WriteString(m.common.Global.AuthState.Cookie)
	sb.WriteString("\n")
	sb.WriteString(m.debug)
	sb.WriteString("\n")

	if m.stage == picker {
		sb.WriteString(m.buttons.View())
	} else {

		if m.stage == signIn {
			sb.WriteString("Sign in")
		} else if m.stage == signUp {
			sb.WriteString("Sign up")
		}

		sb.WriteString("\n")
		sb.WriteString(m.inputs.View())

	}

	// sections = append(sections, m.global.AuthState.Cookie)

	return sb.String()
}

func signUpInputs(c common.Common) []common.Component {
	inputs := make([]common.Component, 0, 5)

	ic := common.Common{
		Width:  c.Width,
		Height: 3, // TODO does nothing
		Global: c.Global,
	}

	for i := 0; i < 4; i++ {
		input := textfield.New(ic)
		input.CursorStyle(cursorStyle)
		input.CharLimit(80)

		switch i {
		case 0:
			input.Placeholder("Name")
			input.Focus()
		case 1:
			input.Placeholder("Email")
		case 2:
			input.Placeholder("Password")
			input.EchoMode(textinput.EchoPassword)
		case 3:
			input.Placeholder("Retype password")
			input.EchoMode(textinput.EchoPassword)
		}

		inputs = append(inputs, input)
	}

	bc := common.Common{
		Width:  c.Width,
		Height: 1, // TODO does nothing
		Global: c.Global,
	}

	button := button.New(bc, "Sign up", func() tea.Msg {
		if inputs[2].(*textfield.Model).Value() != inputs[3].(*textfield.Model).Value() {
			return signUpRes{false, "non mathcing password"}
		}

		return postSignUp(bc.Global.HttpClient, signUpData{
			inputs[0].(*textfield.Model).Value(),
			inputs[1].(*textfield.Model).Value(),
			inputs[2].(*textfield.Model).Value(),
		})
	})

	inputs = append(inputs, button)

	return inputs
}

func signInInputs(c common.Common) []common.Component {
	inputs := make([]common.Component, 0, 3)

	ic := common.Common{
		Width:  c.Width,
		Height: 3, // TODO does nothing
		Global: c.Global,
	}

	for i := 0; i < 2; i++ {
		input := textfield.New(ic)
		input.CursorStyle(cursorStyle)
		input.CharLimit(80)

		switch i {
		case 0:
			input.Placeholder("Email")
			input.Focus()
		case 1:
			input.Placeholder("Password")
			input.EchoMode(textinput.EchoPassword)
		}

		inputs = append(inputs, input)
	}

	bc := common.Common{
		Width:  c.Width,
		Height: 1, // TODO does nothing
		Global: c.Global,
	}

	button := button.New(bc, "Sign in", func() tea.Msg {
		return postSignIn(bc.Global.HttpClient, signInData{
			inputs[0].(*textfield.Model).Value(),
			inputs[1].(*textfield.Model).Value(),
		})
	})

	inputs = append(inputs, button)

	return inputs
}
