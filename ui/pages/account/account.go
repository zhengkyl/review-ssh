package account

import (
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
}

const (
	picker = iota
	signIn
	signUp
)

type signInMsg struct{}
type signUpMsg struct{}

func New(c common.Common) *Model {
	b := vlist.New(c,
		[]common.Component{
			button.New(c, "     Sign in     ", func() tea.Msg { return signInMsg{} }),
			button.New(c, "     Sign up     ", func() tea.Msg { return signUpMsg{} }),
			button.New(c, "Continue as guest", func() tea.Msg { return common.GuestAuthState }),
		},
	)

	b.Style.Active = lipgloss.NewStyle().PaddingLeft(2)

	inputs := make([]common.Component, 3)

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
		}
		inputs[i] = input
	}

	inputs[2] = button.New(c, "Submit", func() tea.Msg {
		return postAuth(&c.Global.HttpClient, loginData{
			inputs[0].(*textfield.Model).Value(),
			inputs[1].(*textfield.Model).Value(),
		})

	})

	m := &Model{
		common:     c,
		inputs:     vlist.New(c, inputs),
		buttons:    b,
		focusIndex: 0,
	}

	return m
}

func (m *Model) SetSize(width, height int) {
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

	case signInMsg:
		m.stage = signIn
	case signUpMsg:
		m.stage = signUp
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
