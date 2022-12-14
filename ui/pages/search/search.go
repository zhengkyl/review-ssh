package search

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type SearchModel struct {
	common common.Common
	input  textinput.Model
}

func New(common common.Common) *SearchModel {

	input := textinput.New()
	input.Placeholder = "Search for movies and shows..."
	input.Focus()
	input.CharLimit = 80

	m := &SearchModel{
		input:  input,
		common: common,
	}

	m.SetSize(common.Width, common.Height)

	return m
}

func (m *SearchModel) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	// wm, hm := m.getMargins()

}

func (m *SearchModel) getMargins() (wm, hm int) {
	wm = 0
	hm = 0

	return
}

func (m *SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m *SearchModel) View() string {
	var view string

	wm, _ := m.getMargins()

	ss := lipgloss.NewStyle().Width(m.common.Width - wm)
	view = ss.Render(m.input.View())

	return view
}

// var cmd tea.Cmd

// m.searchInput, cmd = m.searchInput.Update(msg)
