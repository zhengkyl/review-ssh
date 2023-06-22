package reviewlist

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
	"github.com/zhengkyl/review-ssh/ui/util"
)

type Model struct {
	props        common.Props
	reviews      []common.Review
	offset       int
	active       int
	visibleItems int

	itemSpinner spinner.Model
	spinning    bool
}

var (
	normalStyle = lipgloss.NewStyle()
	activeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	listStyle   = lipgloss.NewStyle().Margin(1)
	dotdotdot   = spinner.Spinner{Frames: []string{"", ".", ".. ", "...", "..", "."}, FPS: time.Second / 3}
)

func New(p common.Props) *Model {
	m := &Model{
		props:       p,
		active:      0,
		itemSpinner: spinner.New(spinner.WithSpinner(dotdotdot)),
		spinning:    false,
	}
	m.SetSize(p.Width, p.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	hf := listStyle.GetHorizontalFrameSize()
	vf := listStyle.GetVerticalFrameSize()

	m.props.Width = width - hf
	m.props.Height = height - vf

	m.visibleItems = util.Max(height/2, 0)

	// Try to keep active item same pos from top when resizing
	maxIndex := util.Max(m.visibleItems-1, 0)
	newIndex := util.Min(m.active-m.offset, maxIndex)
	m.offset = m.active - newIndex
}

func (m *Model) SetReviews(reviews []common.Review) {
	m.reviews = reviews
	m.active = 0
	m.offset = 0
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	if len(m.reviews) == 0 {
		return m, nil
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case spinner.TickMsg:
		if m.spinning {
			m.itemSpinner, cmd = m.itemSpinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		prevActive := m.active
		switch {
		case key.Matches(msg, m.props.Global.KeyMap.Down):
			m.active = util.Min(m.active+1, len(m.reviews)-1)

			if m.active == m.offset+m.visibleItems {
				m.offset++
			}
		case key.Matches(msg, m.props.Global.KeyMap.Up):
			m.active = util.Max(m.active-1, 0)

			if m.active == m.offset-1 {
				m.offset = m.active
			}
		case key.Matches(msg, m.props.Global.KeyMap.Select):
			cmd = func() tea.Msg {
				return common.ShowPage{
					Category: m.reviews[m.active].Category,
					Tmdb_id:  m.reviews[m.active].Tmdb_id,
					Season:   m.reviews[m.active].Season,
				}
			}
			cmds = append(cmds, cmd)
		}

		if prevActive != m.active {

		}
	}

	itemsLoading := false

	for i := m.offset; i < m.offset+m.visibleItems+2 && i < len(m.reviews); i++ {
		review := m.reviews[i]

		switch review.Category {
		case enums.Film:
			ok, loading, _ := m.props.Global.FilmCache.Get(review.Tmdb_id)
			if ok {
				continue
			}

			if loading {
				itemsLoading = true
			} else {
				m.props.Global.FilmCache.SetLoading(review.Tmdb_id)
				cmds = append(cmds, common.GetFilmCmd(m.props.Global, review.Tmdb_id))
			}

		case enums.Show:
			ok, loading, _ := m.props.Global.ShowCache.Get(review.Tmdb_id)
			if ok {
				continue
			}

			if loading {
				itemsLoading = true
			} else {
				m.props.Global.ShowCache.SetLoading(review.Tmdb_id)
				cmds = append(cmds, common.GetShowCmd(m.props.Global, review.Tmdb_id))
			}
		}

	}

	if m.spinning && !itemsLoading {
		m.spinning = false
	} else if !m.spinning && itemsLoading {
		m.spinning = true
		cmds = append(cmds, m.itemSpinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	spinner := m.itemSpinner.View()
	if len(m.reviews) == 0 {
		return "Loading reviews..."
	}

	viewSb := strings.Builder{}

	// 8 wide review
	// 13 status
	// 3 gaps
	// 3 wide scrollbar
	titleWidth := m.props.Width - 8 - 13 - 3 - 3

	for i := m.offset; i < m.offset+m.visibleItems && i < len(m.reviews); i++ {

		review := m.reviews[i]

		title := "Loading" + spinner

		switch review.Category {
		case enums.Film:
			ok, _, film := m.props.Global.FilmCache.Get(review.Tmdb_id)
			if ok {
				title = film.Title
			}
		case enums.Show:
			ok, _, show := m.props.Global.ShowCache.Get(review.Tmdb_id)
			if ok {
				title = show.Name
			}
		}

		sectionSb := strings.Builder{}

		sectionSb.WriteString(util.TruncOrPadASCII(title, titleWidth))

		sectionSb.WriteString(common.RenderThinRating(review.Fun_before, review.Fun_during, review.Fun_after))
		sectionSb.WriteString(" ")

		sectionSb.WriteString(util.TruncOrPadASCII(review.Status.String(), 13))
		sectionSb.WriteString(" ")

		// sectionSb.WriteString(util.TruncOrPadASCII(review.Text, 20))

		sectionSb.WriteString("\n")

		section := sectionSb.String()

		if i == m.active {
			section = activeStyle.Render(section)
		} else {
			section = normalStyle.Render(section)
		}

		if i > m.offset {
			viewSb.WriteString("\n")
		}

		viewSb.WriteString(section)
	}

	scrollPositions := len(m.reviews) - m.visibleItems + 1 // initial + all nonvisible
	scrollBar := renderScrollbar(m.props.Height, scrollPositions, m.offset)

	return lipgloss.JoinHorizontal(lipgloss.Top, listStyle.Render(viewSb.String()), scrollBar)
}
