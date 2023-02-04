package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
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

type UiModel struct {
	common     common.Common
	tabs       []common.Component
	activeTab  int
	httpClient *retryablehttp.Client
}

func New(httpClient *retryablehttp.Client) *UiModel {

	return &UiModel{
		tabs: make([]common.Component, 1),
		common: common.Common{
			// Width: ,
			Styles: styles.DefaultStyles(),
		},
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
	// m.tabs[0] = poster.New(m.common, "https://image.tmdb.org/t/p/w92/kgwjIb2JDHRhNk13lmSxiClFjVk.jpg")

	m.SetSize(m.common.Width, m.common.Height)

	cmds := []tea.Cmd{
		// m.tabs[0].Init(),
		m.tabs[searchPage].Init(),
	}

	return tea.Batch(cmds...)
}

func (m UiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

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

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

			// The "up" and "k" keys move the cursor up
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
