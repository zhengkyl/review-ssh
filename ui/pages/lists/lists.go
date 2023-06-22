package lists

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
	"github.com/zhengkyl/review-ssh/ui/pages/lists/reviewlist"
)

var tabNames = []string{
	"All",
	"Watching",
	"Plan To Watch",
	"Completed",
	"Dropped",
}

// This must match the order of tabNames
var tabStatuses = []enums.Status{
	255, // This should never be accessed
	enums.Watching,
	enums.PlanToWatch,
	enums.Completed,
	enums.Dropped,
}

var NUM_LISTS = len(tabNames)

var (
	tabBorder      = lipgloss.NormalBorder()
	tabStyle       = lipgloss.NewStyle().Padding(0, 1).Border(tabBorder, true)
	activeTabStyle = lipgloss.NewStyle().Padding(0, 1).BorderForeground(lipgloss.Color("#7D56F4")).Border(tabBorder, true)
)

type Model struct {
	props     common.Props
	activeTab int
	list      *reviewlist.Model
	err       string
}

func New(p common.Props) *Model {
	return &Model{
		props:     p,
		activeTab: 0,
		list:      reviewlist.New(p),
	}
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	m.list.SetSize(width, height-3)
}

type Init struct{}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Init:
		return m, getReviewsCmd(m.props.Global.HttpClient, m.props.Global.AuthState.User.Id)
	case common.GetResponse[common.PageResult[common.Review]]:
		if msg.Ok {
			reviews := make([]common.Review, 0, len(m.props.Global.ReviewMap))
			for _, review := range m.props.Global.ReviewMap {
				reviews = append(reviews, review)
			}

			sort.Sort(common.ByUpdatedAt(reviews))
			m.list.SetReviews(reviews)
		} else {
		}
	case tea.KeyMsg:
		switch m.activeTab {
		}
		prevActive := m.activeTab
		if key.Matches(msg, m.props.Global.KeyMap.NextTab) {
			m.activeTab = (m.activeTab + 1) % NUM_LISTS
		} else if key.Matches(msg, m.props.Global.KeyMap.PrevTab) {
			m.activeTab = (m.activeTab - 1 + NUM_LISTS) % NUM_LISTS
		}

		if m.activeTab != prevActive {
			filtered := make([]common.Review, 0)

			if m.activeTab == 0 {
				for _, review := range m.props.Global.ReviewMap {
					filtered = append(filtered, review)
				}
			} else {
				for _, review := range m.props.Global.ReviewMap {
					if tabStatuses[m.activeTab] == review.Status {
						filtered = append(filtered, review)
					}
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

	view.WriteString(tabs)
	view.WriteString("\n")

	view.WriteString(m.list.View())

	view.WriteString(m.err)

	return view.String()
}
