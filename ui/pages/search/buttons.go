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
	buttons      []*button.ButtonModel
}

func NewButtons(common common.Common) *ButtonsModel {
	return &ButtonsModel{
		common: common,
		buttons: []*button.ButtonModel{
			button.New(common, "Show more",
				func() tea.Msg {
					return nil
				},
			),
			button.New(common, "LIKE",
				func() tea.Msg {
					return nil
				},
			),
			button.New(common, "",
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

func (m *ButtonsModel) Update(msg tea.Msg) (*ButtonsModel, tea.Cmd) {
	var cmds []tea.Cmd

	// m.buttons[m.activeButton].text = "no"

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
