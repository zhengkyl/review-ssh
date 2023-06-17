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
	common       common.Common
	reviews      []common.Review
	filmMap      map[int]common.Film
	showMap      map[int]common.Show
	inflight     map[int]struct{}
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
	dotdotdot   = spinner.Spinner{Frames: []string{".", ".. ", "...", ".."}, FPS: time.Second / 3}
)

func New(c common.Common) *Model {
	m := &Model{
		common:      c,
		filmMap:     map[int]common.Film{},
		showMap:     map[int]common.Show{},
		inflight:    map[int]struct{}{},
		active:      0,
		itemSpinner: spinner.New(spinner.WithSpinner(dotdotdot)),
		spinning:    false,
	}
	m.SetSize(c.Width, c.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	vf := listStyle.GetVerticalFrameSize()
	m.visibleItems = util.Max((height-vf)/2, 0)

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

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.offset++
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

	for i := m.offset; i < m.offset+m.visibleItems+2 && i < len(m.reviews); i++ {
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

	loading := len(m.inflight) > 0

	if m.spinning && !loading {
		m.spinning = false
	} else if !m.spinning && loading {
		m.spinning = true
		cmds = append(cmds, m.itemSpinner.Tick)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if len(m.reviews) == 0 {
		return ""
	}

	spinner := m.itemSpinner.View()

	viewSb := strings.Builder{}

	for i := m.offset; i < m.offset+m.visibleItems && i < len(m.reviews); i++ {

		review := m.reviews[i]

		title := "Loading" + spinner

		switch review.Category {
		case enums.Film:
			film, ok := m.filmMap[review.Tmdb_id]
			if ok {
				title = film.Title
			}
		case enums.Show:
			show, ok := m.showMap[review.Tmdb_id]
			if ok {
				title = show.Name
			}
		}

		sectionSb := strings.Builder{}

		sectionSb.WriteString(util.TruncOrPadASCII(title, 30))

		sectionSb.WriteString(" ")
		sectionSb.WriteString(renderThinRating(review.Fun_before, review.Fun_during, review.Fun_after))

		sectionSb.WriteString(" ")
		sectionSb.WriteString(review.Status.String())

		// sectionSb.WriteString(" ")
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
	scrollBar := renderScrollbar(m.common.Height, scrollPositions, m.offset)

	return lipgloss.JoinHorizontal(lipgloss.Top, listStyle.Render(viewSb.String()), scrollBar)
}
