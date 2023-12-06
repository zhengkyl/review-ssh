package filmitem

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
	"golang.org/x/exp/slices"
)

var (
	itemStyle       = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(2).BorderStyle(lipgloss.Border{Left: " "}).BorderLeft(true)
	activeItemStyle = lipgloss.NewStyle().PaddingLeft(1).PaddingRight(2).Foreground(
		lipgloss.Color("#F25D94")).BorderStyle(lipgloss.Border{Left: "┃"}).
		BorderForeground(lipgloss.Color("#F25D94")).BorderLeft(true)

	textStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#fff"))
	subtitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	contentStyle  = lipgloss.NewStyle().MarginLeft(2)
)

// NOTE: Fullwidth spaces are 2 wide
const POSTER_WIDTH = 4 * 2
const POSTER_HEIGHT = 6

type Model struct {
	props   common.Props
	film    common.Film
	poster  *poster.Model
	focused bool
}

func New(p common.Props, film common.Film) *Model {
	m := &Model{
		p,
		film,
		poster.New(
			common.Props{Width: POSTER_WIDTH, Height: POSTER_HEIGHT, Global: p.Global},
			"https://image.tmdb.org/t/p/w200"+film.Poster_path,
		),
		false,
	}

	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Init() tea.Cmd {
	return m.poster.Init()
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *common.KeyEvent:
		switch {
		case key.Matches(msg.KeyMsg, m.props.Global.KeyMap.Select):
			msg.Handled = true
			return m, func() tea.Msg {
				return common.ShowFilm(m.film.Id)
			}
		}
	}

	_, cmd := m.poster.Update(msg)
	return m, cmd
}

func (m *Model) View() string {

	contentWidth := m.props.Width - itemStyle.GetHorizontalFrameSize() - POSTER_WIDTH - contentStyle.GetHorizontalFrameSize()

	// Subtract 15 to account for long word causing early newline.
	desc := ellipsisText(m.film.Overview, contentWidth*2-15)

	var releaseYear string
	if len(m.film.Release_date) > 4 {
		releaseYear = m.film.Release_date[:4]
	}

	str := lipgloss.JoinHorizontal(lipgloss.Top, titleStyle.Render(m.film.Title), " ", subtitleStyle.Render(releaseYear))

	str = lipgloss.JoinVertical(lipgloss.Left, str, textStyle.Width(contentWidth).Render(desc))

	str += "\n\n"

	str = contentStyle.Render(str)

	str = lipgloss.JoinHorizontal(lipgloss.Top, m.poster.View(), str)

	if m.focused {
		// str = lipgloss.JoinHorizontal(lipgloss.Left, "> ", str)
		str = activeItemStyle.Render(str)
	} else {
		str = itemStyle.Render(str)
	}

	return str
}

var ellipsisPos = []rune{' ', '.', ','}

func ellipsisText(s string, max int) string {
	if max >= len(s) {
		return s
	}

	chars := []rune(s)

	// end is an exclusive bound
	var end int
	for end = max - 3; end >= 1; end-- {
		c := chars[end]
		prevC := chars[end-1]

		if slices.Contains(ellipsisPos, c) && !slices.Contains(ellipsisPos, prevC) {
			break
		}
	}

	if end == 0 {
		end = max - 3
	}

	return string(chars[:end]) + "..."
}
