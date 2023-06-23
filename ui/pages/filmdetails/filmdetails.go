package filmdetails

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/dropdown"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
)

var (
	viewStyle = lipgloss.NewStyle().Margin(1)
)

type Model struct {
	props    common.Props
	poster   *poster.Model
	loaded   bool
	filmId   int
	film     common.Film
	dropdown dropdown.Model
}

func New(p common.Props) *Model {
	m := &Model{
		props:  p,
		poster: &poster.Model{},
		loaded: false,
		filmId: 0,
		film:   common.Film{},
		dropdown: *dropdown.New(common.Props{Width: 20, Height: 3, Global: p.Global}, "Add movie", []dropdown.Option{
			{Text: "Plan to Watch", Callback: func() tea.Msg { return nil }},
			{Text: "Completed", Callback: func() tea.Msg { return nil }},
			{Text: "Watching", Callback: func() tea.Msg { return nil }},
			{Text: "Dropped", Callback: func() tea.Msg { return nil }},
		}),
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	hf := viewStyle.GetHorizontalFrameSize()
	vf := viewStyle.GetVerticalFrameSize()

	m.props.Width = width - hf
	m.props.Height = height - vf

	// TODO
	// m.dropdown.SetSize()
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.props.Global.KeyMap.Select):
			m.dropdown.Focus()
		}
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

	if m.dropdown.Focused() {
		_, cmd := m.dropdown.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.loaded {
		return "loading..."
	}

	left := m.poster.View()

	descStyle := lipgloss.NewStyle().Width(m.props.Width - m.poster.Width() - 2).Height(5)

	rightSb := strings.Builder{}
	rightSb.WriteString("\n")
	rightSb.WriteString(m.film.Title + " (" + m.film.Release_date[:4] + ")")
	rightSb.WriteString("\n\n")
	rightSb.WriteString(m.film.Release_date)
	rightSb.WriteString("\n\n")
	rightSb.WriteString(m.dropdown.View())
	rightSb.WriteString("\n\n")

	// see func (r Review) Key() int
	review, ok := m.props.Global.ReviewMap[-m.filmId]

	if ok {
		rightSb.WriteString(common.RenderThickRating(review.Fun_before, review.Fun_during, review.Fun_after))
		rightSb.WriteString("\n\n")
	}

	rightSb.WriteString(descStyle.Render(m.film.Overview))
	rightSb.WriteString("\n\n")

	// return rightSb.String()

	return viewStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", rightSb.String()))
}
