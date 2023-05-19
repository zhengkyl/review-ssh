package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/pages/account"
	"github.com/zhengkyl/review-ssh/ui/pages/lists"
	"github.com/zhengkyl/review-ssh/ui/pages/movie"
	"github.com/zhengkyl/review-ssh/ui/pages/search"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	// highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	docStyle = lipgloss.NewStyle().Background(lipgloss.Color("#7D56F4")).Padding(1, 2)

// windowStyle    = lipgloss.NewStyle().BorderForeground(highlightColor)
)

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fb7185"))
)

type Model struct {
	common      common.Common
	searchField *textfield.Model
	accountPage *account.Model
	listsPage   *lists.Model
	searchPage  *search.Model
	moviePage   *movie.Model
}

func New(c common.Common) *Model {

	searchField := textfield.New(c)
	searchField.CharLimit(80)
	searchField.Placeholder("(s)earch for movies...")

	m := &Model{
		common:      c,
		searchField: searchField,
		accountPage: account.New(c),
		listsPage:   lists.New(c),
		searchPage:  search.New(c),
		moviePage:   movie.New(c),
	}

	m.SetSize(c.Width, c.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	// wm, hm := ui.getMargins()

	m.common.Width = width
	m.common.Height = height

	m.searchField.SetSize(width-20, 3)

	contentHeight := height - 3
	m.accountPage.SetSize(width, contentHeight)
	m.listsPage.SetSize(width, contentHeight)
	m.searchPage.SetSize(width, contentHeight)
	m.moviePage.SetSize(width, contentHeight)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.AuthState:
		m.common.Global.AuthState = msg
	case tea.WindowSizeMsg:
		frameW, frameH := m.common.Global.Styles.App.GetFrameSize()

		viewW, viewH := msg.Width-frameW, msg.Height-frameH

		m.SetSize(viewW, viewH)

	case tea.KeyMsg:
		if key.Matches(msg, m.common.Global.KeyMap.Quit) {
			return m, tea.Quit
		}

		// TODO all other focusables
		if !m.searchField.Focused() {

			if key.Matches(msg, m.common.Global.KeyMap.Search) {
				return m, m.searchField.Focus()
			}

		}

	}

	if m.searchField.Focused() {
		_, cmd := m.searchField.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	view := strings.Builder{}

	bar := lipgloss.JoinHorizontal(lipgloss.Center, titleStyle.Render("movielo"), m.searchField.View())
	view.WriteString(bar)
	view.WriteString("\n")

	if !m.common.Global.AuthState.Authed {
		view.WriteString(m.accountPage.View())
		return view.String()
	}

	parent := m.common.Global.Styles.App.Render(view.String())
	return util.RenderOverlay(parent, docStyle.Render("hello there\nthis should be an overlay\ndid it work?"), 5, 20)

}
