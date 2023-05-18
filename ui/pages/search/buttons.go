package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
)

type ButtonsModel struct {
	common       common.Common
	activeButton int
	buttons      []*button.Model
}

func NewButtons(c common.Common) *ButtonsModel {
	return &ButtonsModel{
		common: c,
		buttons: []*button.Model{
			// button.New(common, "Show more",
			// 	func() tea.Msg {
			// 		return nil
			// 	},
			// ),
			button.New(c, ":D",
				func() tea.Msg {
					return nil
				},
			),
			button.New(c, "<3",
				func() tea.Msg {
					return nil
				},
			),
			button.New(c, "au",
				func() tea.Msg {
					return nil
				},
			),
			button.New(c, "",
				func() tea.Msg {
					return nil
				},
			),
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

func (m *ButtonsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// m.buttons[m.activeButton].text = "no"
	var cmd tea.Cmd
	for _, b := range m.buttons {
		_, cmd = b.Update(msg)

		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *ButtonsModel) View() string {

	var buttons []string

	gap := " "

	for _, button := range m.buttons {
		buttons = append(buttons, button.View(), gap)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, buttons...)
}
