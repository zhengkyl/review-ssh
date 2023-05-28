package account

import (
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
	showButtons = iota
	showInputs
)

func New(c common.Common) *Model {
	b := vlist.New(c,
		[]tea.Model{
			button.New(c, "     Sign in     ", func() tea.Msg { return nil }),
			button.New(c, "     Sign up     ", func() tea.Msg { return nil }),
			button.New(c, "Continue as guest", func() tea.Msg { return nil }),
		},
	)

	b.Style.Active = lipgloss.NewStyle().PaddingLeft(2)

	inputs := make([]tea.Model, 3)

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

	case tea.KeyMsg:
		switch {
		}
	}

	var cmd tea.Cmd

	_, cmd = m.inputs.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = m.buttons.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	// return "-0\n1\n2\n3\n4\n5\n6\n7\n8\n-9"
	// return fmt.Sprintf("ITEM: %p\n1\n2\n3\n4\n5\n6\n7", m)

	// if m.global.AuthState.Authed {
	// 	return m.global.AuthState.User.Name
	// }

	var sections []string

	sections = append(sections, m.common.Global.AuthState.Cookie)

	if m.stage == showInputs {
		sections = append(sections, m.inputs.View())
	}
	if m.stage == showButtons {
		sections = append(sections, m.buttons.View())
	}

	// sections = append(sections, m.global.AuthState.Cookie)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
