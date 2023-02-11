package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Button struct {
	text     string
	callback tea.Cmd
}

type ButtonsModel struct {
	common       common.Common
	activeButton int
	buttons      []Button
}

var (
	buttonStyle       = lipgloss.NewStyle()
	activeButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func NewButtons(common common.Common) *ButtonsModel {
	return &ButtonsModel{
		common: common,
		buttons: []Button{
			{"[ Show more ]", func() tea.Msg {
				return nil
			}},
			{"[ LIKE ]", func() tea.Msg {
				return nil
			}},
			{"[ STAR ]", func() tea.Msg {
				return nil
			}},
			// {"👍", func() tea.Msg {
			// 	return nil
			// }},
			// {"⭐", func() tea.Msg {
			// 	return nil
			// }},
		},
	}
}

func (m *ButtonsModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	return tea.Batch(cmds...)
}

func (m *ButtonsModel) Update(msg tea.Msg) (*ButtonsModel, tea.Cmd) {
	var cmds []tea.Cmd

	// m.buttons[m.activeButton].text = "no"

	return m, tea.Batch(cmds...)
}

func (m *ButtonsModel) View() string {

	var buttons []string

	gap := " "

	for i, button := range m.buttons {
		if i == m.activeButton {
			buttons = append(buttons, activeButtonStyle.Render(button.text), gap)
		} else {
			buttons = append(buttons, buttonStyle.Render(button.text), gap)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
}
