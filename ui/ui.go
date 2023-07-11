package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/button"
	"github.com/zhengkyl/review-ssh/ui/components/dialog"
	"github.com/zhengkyl/review-ssh/ui/pages/account"
	"github.com/zhengkyl/review-ssh/ui/pages/filmdetails"
	"github.com/zhengkyl/review-ssh/ui/pages/lists"
	"github.com/zhengkyl/review-ssh/ui/pages/search"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	appStyle   = lipgloss.NewStyle().MarginBottom(1)
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fb7185")).Padding(0, 1)
	title      = titleStyle.Render("review-ssh")
)

type page int

const (
	ACCOUNT page = iota
	LISTS
	FILMDETAILS
	SEARCH
)

type Model struct {
	props           common.Props
	accountPage     *account.Model
	listsPage       *lists.Model
	filmdetailsPage *filmdetails.Model
	searchPage      *search.Model
	dialog          *dialog.Model
	help            help.Model
	page            page
}

func New(p common.Props) *Model {

	m := &Model{
		props:           p,
		accountPage:     account.New(p),
		listsPage:       lists.New(p),
		filmdetailsPage: filmdetails.New(p),
		searchPage:      search.New(p),
		dialog:          dialog.New(p, "Quit program?"),
		help:            help.New(),
	}

	m.dialog.Buttons(
		*button.New(p, "Yes", tea.Quit),
		*button.New(p, "No", func() tea.Msg {
			m.dialog.Blur()
			return nil
		}))

	m.SetSize(p.Width, p.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	viewW := width
	viewH := height - 2 // 1 for bottom margin + 1 for help

	contentHeight := viewH - 3

	m.accountPage.SetSize(util.Max(viewW/2, 30), contentHeight)

	m.listsPage.SetSize(viewW, contentHeight)
	m.searchPage.SetSize(viewW, contentHeight)
	m.filmdetailsPage.SetSize(viewW, contentHeight)

	m.help.Width = width
}

func (m *Model) Init() tea.Cmd {
	// _, cmd := m.filmdetailsPage.Update(filmdetails.Init(109445))
	// m.page = FILMDETAILS
	// return cmd
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case func():
		msg()
	case common.AuthState:
		m.props.Global.AuthState.Authed = msg.Authed
		m.props.Global.AuthState.Cookie = msg.Cookie
		m.props.Global.AuthState.User = msg.User
		m.page = LISTS
		return m, m.listsPage.Init()
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)

	case common.ShowFilm:
		cmds = append(cmds, m.filmdetailsPage.Init(int(msg)))
		m.page = FILMDETAILS

	case tea.KeyMsg:
		var cmd tea.Cmd
		event := &common.KeyEvent{KeyMsg: msg, Handled: false}

		if m.dialog.Focused() {
			_, cmd = m.dialog.Update(event)
		}

		if event.Handled {
			return m, cmd
		}

		switch m.page {
		case ACCOUNT:
			_, cmd = m.accountPage.Update(event)
		case LISTS:
			_, cmd = m.listsPage.Update(event)
		case FILMDETAILS:
			_, cmd = m.filmdetailsPage.Update(event)
		case SEARCH:
			_, cmd = m.searchPage.Update(msg)
		}

		if event.Handled {
			return m, cmd
		}

		if event.Handled {
			return m, cmd
		}

		if key.Matches(msg, m.props.Global.KeyMap.Quit) {
			if m.dialog.Focused() {
				return m, tea.Quit
			}
			m.dialog.Focus()
		}

		return m, nil
	}

	var cmd tea.Cmd
	_, cmd = m.dialog.Update(msg)
	cmds = append(cmds, cmd)
	switch m.page {
	case ACCOUNT:
		_, cmd = m.accountPage.Update(msg)
	case LISTS:
		_, cmd = m.listsPage.Update(msg)
	case FILMDETAILS:
		_, cmd = m.filmdetailsPage.Update(msg)
	case SEARCH:
		_, cmd = m.searchPage.Update(msg)
	}
	cmds = append(cmds, cmd)

	// m.help, cmd = m.help.Update(msg)
	// cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {

	view := strings.Builder{}
	if !m.props.Global.AuthState.Authed {
		// 3 tall to match search bar + fullwidth to allow centering accountPage view
		rightPad := util.Max(m.props.Width-ansi.PrintableRuneWidth(title), 0)
		appBar := "\n" + title + strings.Repeat(" ", rightPad) + "\n"

		centered := lipgloss.JoinVertical(lipgloss.Center, appBar, m.accountPage.View())
		view.WriteString(centered)
	} else {

		switch m.page {
		case LISTS:
			view.WriteString(m.listsPage.View())
		case FILMDETAILS:
			view.WriteString(m.filmdetailsPage.View())
		}
	}

	vGap := m.props.Height - 2 - lipgloss.Height(view.String())

	if vGap > 0 {
		view.WriteString(strings.Repeat("\n", vGap))
	}

	view.WriteString("\n")
	view.WriteString(m.help.View(m.props.Global.KeyMap))

	app := view.String()

	if m.dialog.Focused() {
		dialogView := m.dialog.View()

		dialogW := lipgloss.Width(dialogView)
		dialogH := lipgloss.Height(dialogView)

		xOffset := util.Max((m.props.Width-dialogW)/2, 0)
		yOffset := util.Max((m.props.Height-dialogH)/2-3, 0)

		app = util.RenderOverlay(app, m.dialog.View(), xOffset, yOffset)
	}

	return appStyle.Render(app)

}
