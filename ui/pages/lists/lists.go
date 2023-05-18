package lists

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

const (
	WATCHING_LIST = iota
	PLAN_TO_WATCH_LIST
	COMPLETED_LIST
	DROPPED_LIST
	NUM_LISTS
)

var tabNames = []string{
	"watching",
	"plan to watch",
	"completed",
	"dropped",
}

var (
	tabBorder      = lipgloss.NormalBorder()
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
	activeTabStyle = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
)

type Model struct {
	common    common.Common
	tabs      []common.Component
	activeTab int
}

func New(c common.Common) *Model {
	return &Model{
		common:    c,
		tabs:      make([]common.Component, NUM_LISTS),
		activeTab: 0,
	}
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	// for _, tab := range m.tabs {
	// 	tab.SetSize(width, height-3)
	// }
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.AuthState:
		m.common.Global.AuthState = msg
	case tea.WindowSizeMsg:
		frameW, frameH := m.common.Global.Styles.App.GetFrameSize()

		viewW, viewH := msg.Width-frameW, msg.Height-frameH

		m.SetSize(viewW, viewH)

		// for _, tab := range m.tabs {
		// 	_, cmd := tab.Update(msg)
		// 	// m.tabs[i] = tabModel.(common.PageComponent)

		// 	tab.SetSize(viewW, viewH-4)

		// 	cmds = append(cmds, cmd)
		// }
	// Is it a key press?
	case tea.KeyMsg:
		switch m.activeTab {
		}

		if key.Matches(msg, m.common.Global.KeyMap.NextTab) {
			m.activeTab = (m.activeTab + 1) % NUM_LISTS
			// return m, nil
		} else if key.Matches(msg, m.common.Global.KeyMap.PrevTab) {
			m.activeTab = (m.activeTab - 1 + NUM_LISTS) % NUM_LISTS
			// return m, nil
		}
	}

	_, cmd := m.tabs[m.activeTab].Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	view := strings.Builder{}
	names := []string{}

	for i, name := range tabNames {
		if i == m.activeTab {
			names = append(names, activeTabStyle.Render(name))
		} else {
			names = append(names, tabStyle.Render(name))
		}
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top,
		names...,
	)

	view.WriteString(tabs + "\n\n")

	view.WriteString(m.tabs[m.activeTab].View())
	return ""
}
