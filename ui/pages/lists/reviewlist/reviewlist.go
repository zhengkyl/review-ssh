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
	common       common.Common
	reviews      []common.Review
	filmMap      map[int]common.Film
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
		filmMap:      map[int]common.Film{},
		inflight:     map[int]struct{}{},
		offset:       0,
		active:       0,
		visibleItems: c.Height / 2,
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height
	m.visibleItems = m.common.Height / 2
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
	case common.Film:
		m.filmMap[msg.Id] = msg
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
		_, ok := m.filmMap[review.Tmdb_id]
		if ok {
			continue
		}
		_, inflight := m.inflight[review.Tmdb_id]
		if inflight {
			continue
		}

		m.inflight[review.Tmdb_id] = struct{}{}

		cmds = append(cmds, func() tea.Msg {
			return getFilm(m.common.Global.HttpClient, m.common.Global.Config.TMDB_API_KEY, review.Tmdb_id)
		})
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	viewSb := strings.Builder{}

	for i := m.offset; i < len(m.reviews); i++ {

		review := m.reviews[i]
		var film common.Film
		film, ok := m.filmMap[review.Tmdb_id]
		if !ok {
			film = loadingFilm
		}

		sectionSb := strings.Builder{}

		sectionSb.WriteString(util.TruncOrPadASCII(film.Title, m.common.Width-50))

		ratingIndex := 0
		if review.Fun_before {
			ratingIndex += 1
		}
		if review.Fun_during {
			ratingIndex += 2
		}
		if review.Fun_after {
			ratingIndex += 4
		}
		sectionSb.WriteString(" ")
		sectionSb.WriteString(ratings[ratingIndex])

		sectionSb.WriteString(" ")
		sectionSb.WriteString(review.Status)
		sectionSb.WriteString("\n")

		section := sectionSb.String()

		if i == m.active {
			section = active.Render(section)
		} else {
			section = normal.Render(section)
		}

		if i > m.offset {
			viewSb.WriteString("\n")
		}

		viewSb.WriteString(section)
	}

	// TODO paginatation

	return viewSb.String()
}

var loadingFilm = common.Film{
	Id:           -1,
	Title:        "Loading",
	Overview:     "Loading description",
	Poster_path:  "poster path? I barely know her",
	Release_date: "0000-00-00",
}
