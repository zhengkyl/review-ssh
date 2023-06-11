package lists

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/pages/lists/reviewlist"
)

const (
	ALL_LIST = iota
	WATCHING_LIST
	PLAN_TO_WATCH_LIST
	COMPLETED_LIST
	DROPPED_LIST
	NUM_LISTS
)

var tabNames = []string{
	"all",
	"watching",
	"plan to watch",
	"completed",
	"dropped",
}

var (
	tabBorder      = lipgloss.NormalBorder()
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Border(tabBorder, true)
	activeTabStyle = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
)

type Model struct {
	common    common.Common
	activeTab int
	list      *reviewlist.Model
}

func New(c common.Common) *Model {
	return &Model{
		common:    c,
		list:      reviewlist.New(c),
		activeTab: 0,
	}
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	m.list.SetSize(width, height-3)
}

func (m *Model) Init() tea.Cmd {
	return getReviewsCmd(m.common.Global.HttpClient, m.common.Global.AuthState.User.Id)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []common.Review:
		m.list.SetReviews(msg)
		return m, nil

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

	_, cmd := m.list.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	view := strings.Builder{}
	names := []string{}

	for i, tabName := range tabNames {
		var name string
		if i == m.activeTab {
			name = activeTabStyle.Render(tabName)
		} else {
			name = tabStyle.Render(tabName)
		}
		names = append(names, name)
	}

	tabs := lipgloss.JoinHorizontal(lipgloss.Top,
		names...,
	)

	view.WriteString(tabs + "\n\n")

	view.WriteString(m.list.View())
	return view.String()
}
