package ui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/keymap"
	"github.com/zhengkyl/review-ssh/ui/pages/account"
	"github.com/zhengkyl/review-ssh/ui/pages/search"
	"github.com/zhengkyl/review-ssh/ui/styles"
)

// var (
// 	testStyle = lipgloss.NewStyle().
// 			Bold(true).
// 			Foreground(lipgloss.Color("#FAFAFA")).
// 			Background(lipgloss.Color("#7D56F4")).
// 			PaddingTop(2).
// 			PaddingLeft(4).
// 			Width(22)
// 	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
// 	docStyle       = lipgloss.NewStyle().Background(lipgloss.Color("#7D56F4")).Padding(1, 2)
// 	windowStyle    = lipgloss.NewStyle().BorderForeground(highlightColor)
// )

const (
	searchPage int = iota
	accountPage
)

const NUM_PAGES = 2

type UiModel struct {
	common     common.Common
	tabs       []common.Component
	activeTab  int
	httpClient *retryablehttp.Client
}

func New(httpClient *retryablehttp.Client) *UiModel {

	return &UiModel{
		common: common.Common{
			// Width: ,
			Styles: styles.DefaultStyles(),
			KeyMap: keymap.DefaultKeyMap(),
		},
		tabs:       make([]common.Component, NUM_PAGES),
		activeTab:  0,
		httpClient: httpClient,
	}
}

func (m *UiModel) SetSize(width, height int) {
	m.common.SetSize(width, height)

	// wm, hm := ui.getMargins()

	// SetSize(width - wm, height - hm)

}

func (m UiModel) Init() tea.Cmd {

	m.tabs[searchPage] = search.New(m.common, m.httpClient)
	m.tabs[accountPage] = account.NewLogin(m.common, m.httpClient)

	m.SetSize(m.common.Width, m.common.Height)

	cmds := []tea.Cmd{
		m.tabs[searchPage].Init(),
		m.tabs[accountPage].Init(),
	}

	return tea.Batch(cmds...)
}

func (m UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		frameW, frameH := m.common.Styles.App.GetFrameSize()

		viewW, viewH := msg.Width-frameW, msg.Height-frameH

		m.SetSize(viewW, viewH)

		for i, t := range m.tabs {
			tabModel, cmd := t.Update(msg)
			m.tabs[i] = tabModel.(common.Component)

			m.tabs[i].SetSize(viewW, viewH)

			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	// Is it a key press?
	case tea.KeyMsg:
		switch m.activeTab {
		case searchPage:
		case accountPage:
		}

		if key.Matches(msg, m.common.KeyMap.NextTab) {
			m.activeTab = (m.activeTab + 1) % NUM_PAGES
		} else if key.Matches(msg, m.common.KeyMap.PrevTab) {
			m.activeTab = (m.activeTab - 1 + NUM_PAGES) % NUM_PAGES
		} else if key.Matches(msg, m.common.KeyMap.Quit) {
			return m, tea.Quit
		}

	}

	tabModel, cmd := m.tabs[m.activeTab].Update(msg)
	m.tabs[m.activeTab] = tabModel.(common.Component)
	if cmd != nil {
		cmds = append(cmds, cmd)
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, tea.Batch(cmds...)
}

func (m UiModel) View() string {

	var view string

	// for i, tab := range m.tabs {

	// }

	// view.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabs[m.activeTab]))
	// The footer

	// view = ui.
	view = m.tabs[m.activeTab].View()
	// view = lipgloss.JoinVertical(lipgloss.Left, ui.)
	// Send the UI for rendering
	return m.common.Styles.App.Render(view)
}
