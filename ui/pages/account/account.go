package account

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("197"))
	cursorStyle  = focusedStyle.Copy()
	accountStyle = lipgloss.NewStyle().Padding(1, 3).Border(lipgloss.RoundedBorder(), true)
)

const SUBMIT_INDEX = 2

type Model struct {
	props      common.Props
	inputs     *vlist.Model
	buttons    *vlist.Model
	focusIndex int
	stage      int
	err        string
	help       help.Model
}

const (
	picker = 0
	signIn = 2
	signUp = 3
)

type signInMsg struct{}
type signUpMsg struct{}

func New(p common.Props) *Model {
	b := vlist.New(p,
		button.New(p, "     Sign in     ", func() tea.Msg { return signInMsg{} }),
		button.New(p, "     Sign up     ", func() tea.Msg { return signUpMsg{} }),
		button.New(p, "Continue as guest", func() tea.Msg { return common.GuestAuthState }),
	)

	b.Style.Active = lipgloss.NewStyle().Margin(1, 0)
	b.Style.Normal = lipgloss.NewStyle().Margin(1, 0)

	m := &Model{
		props:      p,
		inputs:     vlist.New(p),
		buttons:    b,
		focusIndex: 0,
		help:       help.New(),
	}

	return m
}

func (m *Model) Height() int {
	return m.props.Height
}

func (m *Model) Width() int {
	return m.props.Width
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	m.help.Width = width

	m.inputs.SetSize(width, height)
	m.buttons.SetSize(width, height)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case signUpRes:
		m.err = msg.err
	case signInRes:
		m.err = msg.err
	case signInMsg:
		m.stage = signIn
		m.inputs.Children = signInInputs(m.props)
		m.inputs.Active = 0
		m.err = ""
	case signUpMsg:
		m.stage = signUp
		m.inputs.Children = signUpInputs(m.props)
		m.inputs.Active = 0
		m.err = ""
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.props.Global.KeyMap.Back):
			if m.stage != 0 {
				m.stage = 0
			}
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
	sb := strings.Builder{}

	// sb.WriteString(m.props.Global.AuthState.Cookie)
	// sb.WriteString("\n")

	if m.stage == picker {
		sb.WriteString("Ahoy there!\n")
		sb.WriteString(m.buttons.View())
	} else {

		if m.stage == signIn {
			sb.WriteString(" Sign in")
		} else if m.stage == signUp {
			sb.WriteString(" Sign up")
		}
		sb.WriteString("\n\n")

		sb.WriteString(m.inputs.View())
		// if m.err != "" {
		sb.WriteString(errStyle.Render(m.err))
		// }

	}

	return accountStyle.Render(sb.String())
}

func signUpInputs(p common.Props) []common.Component {
	inputs := make([]common.Component, 0, 5)

	ic := common.Props{
		Width:  p.Width,
		Height: 3, // TODO does nothing
		Global: p.Global,
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

	bp := common.Props{
		Width:  p.Width,
		Height: 1, // TODO does nothing
		Global: p.Global,
	}

	button := button.New(bp, "Sign up", func() tea.Msg {
		if inputs[2].(*textfield.Model).Value() != inputs[3].(*textfield.Model).Value() {
			return signUpRes{false, "Passwords do not match."}
		}

		return postSignUp(bp.Global.HttpClient, signUpData{
			inputs[0].(*textfield.Model).Value(),
			inputs[1].(*textfield.Model).Value(),
			inputs[2].(*textfield.Model).Value(),
		})
	})
	button.Style.Normal.Margin(1).MarginBottom(0)
	button.Style.Active.Margin(1).MarginBottom(0)

	inputs = append(inputs, button)

	return inputs
}

func signInInputs(p common.Props) []common.Component {
	inputs := make([]common.Component, 0, 3)

	ic := common.Props{
		Width:  p.Width,
		Height: 3, // TODO does nothing
		Global: p.Global,
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

	bp := common.Props{
		Width:  p.Width,
		Height: 1, // TODO does nothing
		Global: p.Global,
	}

	button := button.New(bp, "Sign in", func() tea.Msg {
		return postSignIn(p.Global.HttpClient, signInData{
			inputs[0].(*textfield.Model).Value(),
			inputs[1].(*textfield.Model).Value(),
		})
	})
	button.Style.Normal.Margin(1).MarginBottom(0)
	button.Style.Active.Margin(1).MarginBottom(0)

	inputs = append(inputs, button)

	return inputs
}
