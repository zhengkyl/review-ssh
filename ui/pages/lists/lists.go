package lists

import (
	"sort"
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
	"All",
	"Watching",
	"Plan to Watch",
	"Completed",
	"Dropped",
}

var (
	tabBorder      = lipgloss.NormalBorder()
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Border(tabBorder, true)
	activeTabStyle = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
)

type Model struct {
	common        common.Common
	activeTab     int
	filmReviewMap map[int]common.Review
	list          *reviewlist.Model
	// showReviewMap map[int]common.Review
}

func New(c common.Common) *Model {
	return &Model{
		common:        c,
		activeTab:     0,
		filmReviewMap: make(map[int]common.Review),
		list:          reviewlist.New(c),
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
		for _, review := range msg {
			m.filmReviewMap[review.Tmdb_id] = review
		}
		m.list.SetReviews(msg)

	case tea.KeyMsg:
		switch m.activeTab {
		}
		prevActive := m.activeTab
		if key.Matches(msg, m.common.Global.KeyMap.NextTab) {
			m.activeTab = (m.activeTab + 1) % NUM_LISTS
		} else if key.Matches(msg, m.common.Global.KeyMap.PrevTab) {
			m.activeTab = (m.activeTab - 1 + NUM_LISTS) % NUM_LISTS
		}

		if m.activeTab != prevActive {
			filtered := make([]common.Review, 0)
			for _, review := range m.filmReviewMap {
				if m.activeTab == 0 || review.Status == tabNames[m.activeTab] {
					filtered = append(filtered, review)
				}
			}

			sort.Sort(common.ByUpdatedAt(filtered))

			m.list.SetReviews(filtered)
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
