package filmdetails

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/checkbox"
	"github.com/zhengkyl/review-ssh/ui/components/dropdown"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	viewStyle = lipgloss.NewStyle().Margin(1)
)

type Model struct {
	props       common.Props
	poster      *poster.Model
	loaded      bool
	filmId      int
	film        common.Film
	inputs      []common.Focusable
	dropdown    *dropdown.Model
	checkBefore *checkbox.Model
	checkDuring *checkbox.Model
	checkAfter  *checkbox.Model
	focusIndex  int
}

func New(p common.Props) *Model {
	m := &Model{
		props:  p,
		poster: &poster.Model{},
		loaded: false,
		filmId: 0,
		film:   common.Film{},
		dropdown: dropdown.New(common.Props{Width: 20, Height: 3, Global: p.Global}, "Add movie", []dropdown.Option{
			{Text: "Plan to Watch", Callback: func() tea.Msg { return nil }},
			{Text: "Completed", Callback: func() tea.Msg { return nil }},
		}),
		checkBefore: checkbox.New(p),
		checkDuring: checkbox.New(p),
		checkAfter:  checkbox.New(p),
		inputs:      []common.Focusable{},
		focusIndex:  0,
	}
	m.inputs = append(m.inputs, m.dropdown, m.checkBefore, m.checkDuring, m.checkAfter)

	return m
}

func (m *Model) SetSize(width, height int) {
	hf := viewStyle.GetHorizontalFrameSize()
	vf := viewStyle.GetVerticalFrameSize()

	m.props.Width = width - hf
	m.props.Height = height - vf
}

func (m *Model) Init(filmId int) tea.Cmd {
	m.filmId = filmId
	m.loaded = false
	review, ok := m.props.Global.ReviewMap[m.filmId]
	m.checkBefore.Checked = review.Fun_before
	m.checkDuring.Checked = review.Fun_during
	m.checkAfter.Checked = review.Fun_after

	m.dropdown.Focus()

	if ok {
		return nil
	}
	return func() tea.Msg {
		res := common.GetMyFilmReviewCmd(m.props.Global, m.filmId)().(common.GetResponse[common.PageResult[common.Review]])
		if res.Ok && len(res.Data.Results) > 0 {
			review := res.Data.Results[0]
			m.checkBefore.Checked = review.Fun_before
			m.checkDuring.Checked = review.Fun_during
			m.checkAfter.Checked = review.Fun_after

			m.dropdown.Focus()
		}
		return res
	}
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case *common.KeyEvent:
		prevFocus := m.focusIndex
		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Right):
			m.focusIndex = util.Min(m.focusIndex+1, 3)
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Left):
			m.focusIndex = util.Max(m.focusIndex-1, 0)
		}
		if m.focusIndex != prevFocus {
			m.inputs[m.focusIndex].Focus()
			m.inputs[prevFocus].Blur()
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

	_, cmd := m.inputs[m.focusIndex].Update(msg)
	cmds = append(cmds, cmd)

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

	// make place holder for dropdown which needs to be overlaid
	dropdownView := m.dropdown.View()
	inputs := lipgloss.JoinHorizontal(lipgloss.Top, strings.Repeat(" ", lipgloss.Width(dropdownView)), " ", m.checkBefore.View(), " ", m.checkDuring.View(), " ", m.checkAfter.View())
	rightSb.WriteString(inputs)

	rightSb.WriteString("\n\n")

	rightSb.WriteString(descStyle.Render(m.film.Overview))
	rightSb.WriteString("\n\n")

	rightView := util.RenderOverlay(rightSb.String(), dropdownView, 0, 3)

	// return rightSb.String()

	return viewStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", rightView))
}
