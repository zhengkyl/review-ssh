package filmdetails

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
)

type Model struct {
	props  common.Props
	review common.Review
	poster *poster.Model
	loaded bool
	filmId int
	film   common.Film
}

func New(p common.Props) *Model {
	m := &Model{
		props:  p,
		loaded: false,
	}

	return m
}

func (m *Model) SetFilm(filmId int) {
	m.filmId = filmId
	m.loaded = false
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height
}

func (m *Model) Init() tea.Cmd {
	return m.poster.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if !m.loaded {

		ok, loading, film := m.props.Global.FilmCache.Get(m.filmId)

		if !ok {
			cmd := common.GetFilmCmd(m.props.Global, m.filmId)
			cmds = append(cmds, cmd)
		} else if loading {

		} else {
			m.loaded = true
			m.film = film
			m.poster = poster.New(common.Props{Width: 16, Height: 12, Global: m.props.Global}, "https://image.tmdb.org/t/p/w200"+film.Poster_path)
			cmds = append(cmds, m.poster.Init())

			for _, review := range m.props.Global.ReviewMap {
				if review.Category == enums.Film && review.Tmdb_id == m.filmId {
					m.review = review
					break
				}
			}
		}
	} else {
		_, cmd := m.poster.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.loaded {
		return m.poster.View()
	}
	return "loading..."
}
