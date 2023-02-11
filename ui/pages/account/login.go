package account

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type LoginModel struct {
	common        common.Common
	httpClient    *retryablehttp.Client
	emailInput    textinput.Model
	passwordInput textinput.Model
}

func NewLogin(common common.Common, httpClient *retryablehttp.Client) *LoginModel {
	emailInput := textinput.New()
	emailInput.Placeholder = "Email"
	emailInput.Focus()
	emailInput.CharLimit = 80

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.CharLimit = 80

	return &LoginModel{
		common,
		httpClient,
		emailInput,
		passwordInput,
	}
}

func (m *LoginModel) SetSize(width, height int) {
}

func (m *LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg.(type) {
	case tea.WindowSizeMsg:

	case tea.KeyMsg:
	}

	var cmd tea.Cmd
	m.emailInput, cmd = m.emailInput.Update(msg)
	cmds = append(cmds, cmd)

	m.passwordInput, cmd = m.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *LoginModel) View() string {
	var sections []string

	sections = append(sections, m.emailInput.View(), m.passwordInput.View())

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}
