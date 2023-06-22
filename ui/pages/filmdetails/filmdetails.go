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
	props  common.Props
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

func (m *Model) SetSize(width, height int) {
	hf := viewStyle.GetHorizontalFrameSize()
	vf := viewStyle.GetVerticalFrameSize()

	m.props.Width = width - hf
	m.props.Height = height - vf
}

type Init int

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case Init:
		m.filmId = int(msg)
		m.loaded = false
		_, ok := m.props.Global.ReviewMap[-m.filmId]
		if ok {
			return m, nil
		}
		return m, common.GetMyFilmReviewCmd(m.props.Global, m.filmId)
	}

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

			_, cmd := m.poster.Update(poster.Init{})
			cmds = append(cmds, cmd)

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

	// see func (r Review) Key() int
	review, ok := m.props.Global.ReviewMap[-m.filmId]

	if ok {
		rightSb.WriteString(common.RenderThickRating(review.Fun_before, review.Fun_during, review.Fun_after))
		rightSb.WriteString("\n\n")
	}

	// return rightSb.String()

	return viewStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, " ", rightSb.String()))
}
