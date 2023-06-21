package filmdetails

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
)

var (
	viewStyle = lipgloss.NewStyle().Margin(1)
)

type Model struct {
	props common.Props
	// review common.Review
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
	hf := viewStyle.GetHorizontalFrameSize()
	vf := viewStyle.GetVerticalFrameSize()

	m.props.Width = width - hf
	m.props.Height = height - vf
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
			m.poster = poster.New(common.Props{Width: 28, Height: 21, Global: m.props.Global}, "https://image.tmdb.org/t/p/w200"+film.Poster_path)
			cmds = append(cmds, m.poster.Init())

			// for _, review := range m.props.Global.ReviewMap {
			// 	if review.Category == enums.Film && review.Tmdb_id == m.filmId {
			// 		m.review = review
			// 		break
			// 	}
			// }
		}
	} else {
		_, cmd := m.poster.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.loaded {
		return "loading..."
	}

	left := m.poster.View()

	descStyle := lipgloss.NewStyle().Width(m.props.Width - 16 - 1).Height(5)

	rightSb := strings.Builder{}
	rightSb.WriteString(m.film.Title + " (" + m.film.Release_date[:4] + ")")
	rightSb.WriteString("\n\n")
	rightSb.WriteString(m.film.Release_date)
	rightSb.WriteString("\n\n")
	rightSb.WriteString(descStyle.Render(m.film.Overview))
	rightSb.WriteString("\n\n")

	// return rightSb.String()

	return viewStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, " ", rightSb.String()))
}
