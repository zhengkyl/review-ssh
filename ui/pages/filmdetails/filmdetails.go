package filmdetails

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
	"github.com/zhengkyl/review-ssh/ui/components/checkbox"
	"github.com/zhengkyl/review-ssh/ui/components/dropdown"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	viewStyle = lipgloss.NewStyle().Margin(1)
)

type Model struct {
	props        common.Props
	poster       *poster.Model
	filmLoaded   bool
	reviewLoaded bool
	filmId       int
	inputs       []common.Focusable
	dropdown     *dropdown.Model
	checkBefore  *checkbox.Model
	checkDuring  *checkbox.Model
	checkAfter   *checkbox.Model
	focusIndex   int
	updates      map[string]string
}

func New(p common.Props) *Model {
	m := &Model{
		props:      p,
		poster:     &poster.Model{},
		filmLoaded: false,
		filmId:     0,
		dropdown: dropdown.New(common.Props{Width: 20, Height: 3, Global: p.Global}, "Add movie", []dropdown.Option{
			{Text: "Plan to Watch", Callback: func() tea.Msg { return nil }},
			{Text: "Completed", Callback: func() tea.Msg { return nil }},
		}),
		checkBefore: checkbox.New(p),
		checkDuring: checkbox.New(p),
		checkAfter:  checkbox.New(p),
		inputs:      []common.Focusable{},
		focusIndex:  0,
		updates:     make(map[string]string),
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

func (m *Model) updateInputs(review common.Review) {
	m.checkBefore.Checked = review.Fun_before
	m.checkDuring.Checked = review.Fun_during
	m.checkAfter.Checked = review.Fun_after

	switch review.Status {
	case enums.PlanToWatch:
		m.dropdown.Selected = 0
	case enums.Completed:
		m.dropdown.Selected = 1
	}
}

func (m *Model) Init(filmId int) tea.Cmd {
	m.filmId = filmId
	m.filmLoaded = false
	review, ok := m.props.Global.ReviewMap[m.filmId]
	m.updateInputs(review)

	m.dropdown.Focus()

	if ok {
		m.reviewLoaded = true
		return nil
	}
	return common.GetMyFilmReviewCmd(m.props.Global, m.filmId, func(data common.Paged[common.Review], err error) tea.Msg {
		if err == nil && len(data.Results) > 0 {
			review := data.Results[0]
			m.props.Global.ReviewMap[review.Tmdb_id] = review
			m.updateInputs(review)
		}
		return nil
	})
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

	if !m.filmLoaded {

		ok, loading, film := m.props.Global.FilmCache.Get(m.filmId)

		if !ok {
			cmd := common.GetFilmCmd(m.props.Global, m.filmId)
			cmds = append(cmds, cmd)
		} else if loading {

		} else {
			m.filmLoaded = true
			m.poster = poster.New(common.Props{Width: 28, Height: 21, Global: m.props.Global}, "https://image.tmdb.org/t/p/w200"+film.Poster_path)

			cmds = append(cmds, m.poster.Init())

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
	if !m.filmLoaded {
		return "loading..."
	}
	_, _, film := m.props.Global.FilmCache.Get(m.filmId)

	left := m.poster.View()

	descStyle := lipgloss.NewStyle().Width(m.props.Width - 28 - 2).Height(5)

	rightSb := strings.Builder{}
	rightSb.WriteString("\n")
	rightSb.WriteString(film.Title + " (" + film.Release_date[:4] + ")")
	rightSb.WriteString("\n\n")

	// make place holder for dropdown which needs to be overlaid
	dropdownView := m.dropdown.View()
	inputs := lipgloss.JoinHorizontal(lipgloss.Top, strings.Repeat(" ", lipgloss.Width(dropdownView)), " ", m.checkBefore.View(), " ", m.checkDuring.View(), " ", m.checkAfter.View())
	rightSb.WriteString(inputs)

	rightSb.WriteString("\n\n")

	rightSb.WriteString(descStyle.Render(film.Overview))
	rightSb.WriteString("\n\n")

	rightView := util.RenderOverlay(rightSb.String(), dropdownView, 0, 3)

	return viewStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", rightView))
}
