package reviewlist

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	common    common.Common
	reviews   []common.Review
	movieInfo map[int]common.Movie
	inflight  map[int]struct{}
	// showInfo     map[int]common.Show
	Style        Style
	offset       int
	active       int
	visibleItems int
}

type Style struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

func New(c common.Common) *Model {
	m := &Model{
		common:    c,
		movieInfo: map[int]common.Movie{},
		inflight:  map[int]struct{}{},
		Style: Style{
			Normal: lipgloss.NewStyle(),
			Active: lipgloss.NewStyle(),
		},
		offset:       0,
		active:       0,
		visibleItems: c.Height, // set to highest possible ie 1 height items, set in Update()
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
}

func (m *Model) SetReviews(reviews []common.Review) {
	m.reviews = reviews
	m.active = 0
	m.offset = 0
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.reviews) == 0 {
		return m, nil
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {

	case common.Movie:
		m.movieInfo[msg.Id] = msg

	// case common.Show:
	// 	m.showInfo[msg.Id] = msg

	case tea.KeyMsg:
		prevActive := m.active
		switch {
		case key.Matches(msg, m.common.Global.KeyMap.Down):
			m.active = util.Min(m.active+1, len(m.reviews)-1)

			if m.active == m.offset+m.visibleItems {
				m.offset = m.active
			}
		case key.Matches(msg, m.common.Global.KeyMap.Up):
			m.active = util.Max(m.active-1, 0)

			if m.active == m.offset-1 {
				m.offset = m.active
			}
		}

		if prevActive != m.active {

		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}
	height := 1

	m.visibleItems = 0

	for i := m.offset; i < len(m.reviews); i++ {

		section := m.renderReview(m.reviews[i])

		if i == m.active {
			section = m.Style.Active.Render(section)
		} else {
			section = m.Style.Normal.Render(section)
		}

		sectionHeight := lipgloss.Height(section)

		if height+sectionHeight > m.common.Height {
			break
		}

		height += sectionHeight
		m.visibleItems++

		if i > m.offset {
			sb.WriteString("\n")
		}

		sb.WriteString(section)
	}

	// TODO paginatation
	sb.WriteString("")

	return sb.String()
}

func (m *Model) renderReview(review common.Review) string {
	return ""
}
