package reviewlist

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	common       common.Common
	reviews      []common.Review
	filmMap      map[int]common.Film
	showMap      map[int]common.Show
	inflight     map[int]struct{}
	offset       int
	active       int
	visibleItems int
}

var (
	normal = lipgloss.NewStyle()
	active = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func New(c common.Common) *Model {
	m := &Model{
		common:       c,
		filmMap:      map[int]common.Film{},
		showMap:      map[int]common.Show{},
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
	case common.GetResponse[common.Film]:
		if msg.Ok {
			film := msg.Data
			m.filmMap[film.Id] = film
			delete(m.inflight, film.Id)
		} else {
		}
	case common.GetResponse[common.Show]:
		if msg.Ok {
			show := msg.Data
			m.showMap[show.Id] = show
			delete(m.inflight, show.Id)
		} else {
		}
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

		switch review.Category {
		case enums.Film:
			_, ok := m.filmMap[review.Tmdb_id]
			if ok {
				continue
			}

			_, inflight := m.inflight[review.Key()]
			if inflight {
				continue
			}
			m.inflight[review.Key()] = struct{}{}

			cmds = append(cmds, getFilmCmd(m.common.Global, review.Tmdb_id))

		case enums.Show:
			_, ok := m.showMap[review.Tmdb_id]
			if ok {
				continue
			}

			_, inflight := m.inflight[review.Key()]
			if inflight {
				continue
			}
			m.inflight[review.Key()] = struct{}{}

			cmds = append(cmds, getShowCmd(m.common.Global, review.Tmdb_id))
		}

	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	viewSb := strings.Builder{}

	for i := m.offset; i < len(m.reviews); i++ {

		review := m.reviews[i]

		info := loadingMedia

		switch review.Category {
		case enums.Film:
			film, ok := m.filmMap[review.Tmdb_id]
			if ok {
				info.Title = film.Title
				info.Overview = film.Overview
			}
		case enums.Show:
			show, ok := m.showMap[review.Tmdb_id]
			if ok {
				info.Title = show.Name
				info.Overview = show.Overview
			}
		}

		sectionSb := strings.Builder{}

		sectionSb.WriteString(util.TruncOrPadASCII(info.Title, m.common.Width-50))

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
		sectionSb.WriteString(review.Status.String())
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

var loadingMedia = mediaInfo{
	Title:    "Loading",
	Overview: "Loading description",
}

type mediaInfo struct {
	Title    string
	Overview string
}
