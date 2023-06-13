package reviewlist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	common       common.Common
	reviews      []common.Review
	movieMap     map[int]common.Movie
	inflight     map[int]struct{}
	offset       int
	active       int
	visibleItems int
	results      []res
}

var (
	normal = lipgloss.NewStyle()
	active = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func New(c common.Common) *Model {
	m := &Model{
		common:       c,
		movieMap:     map[int]common.Movie{},
		inflight:     map[int]struct{}{},
		offset:       0,
		active:       0,
		visibleItems: c.Height / 4,
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
	m.visibleItems = m.common.Height / 4
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
	case res:
		m.results = append(m.results, msg)
	case common.Movie:
		m.movieMap[msg.Id] = msg
		delete(m.inflight, msg.Id)
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

	for i := m.offset; i < m.offset+m.visibleItems && i < len(m.reviews); i++ {
		review := m.reviews[i]
		_, ok := m.movieMap[review.Tmdb_id]
		if ok {
			continue
		}
		_, inflight := m.inflight[review.Tmdb_id]
		if inflight {
			continue
		}

		m.inflight[review.Tmdb_id] = struct{}{}

		cmds = append(cmds, func() tea.Msg {
			return getMovie(m.common.Global.HttpClient, m.common.Global.Config.TMDB_API_KEY, review.Tmdb_id)
		})
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	sb := strings.Builder{}
	height := 1

	for i := m.offset; i < len(m.reviews); i++ {

		section := m.renderReview(m.reviews[i])

		if i == m.active {
			section = active.Render(section)
		} else {
			section = normal.Render(section)
		}

		sectionHeight := lipgloss.Height(section)

		if height+sectionHeight > m.common.Height {
			break
		}

		height += sectionHeight

		if i > m.offset {
			sb.WriteString("\n")
		}

		sb.WriteString(section)
	}

	// TODO paginatation
	sb.WriteString("")
	sb.WriteString(fmt.Sprint(m.results))

	return sb.String()
}

var loadingMovie = common.Movie{
	Id:           -1,
	Title:        "Loading",
	Overview:     "Loading description",
	Poster_path:  "poster path? I barely know her",
	Release_date: "0000-00-00",
}

func (m *Model) renderReview(review common.Review) string {
	var movie common.Movie
	movie, ok := m.movieMap[review.Tmdb_id]
	if !ok {
		movie = loadingMovie
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, movie.Title, review.Status, RenderRating(review.Fun_before, review.Fun_during, review.Fun_after))
}
