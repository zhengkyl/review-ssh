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
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/pages/account"
	"github.com/zhengkyl/review-ssh/ui/pages/film"
	"github.com/zhengkyl/review-ssh/ui/pages/lists"
	"github.com/zhengkyl/review-ssh/ui/pages/search"
	"github.com/zhengkyl/review-ssh/ui/util"
)

var (
	titleStyle = lipgloss.NewStyle().Background(lipgloss.Color("#fb7185")).Padding(0, 1)
	title      = titleStyle.Render("review-ssh")
)

type Model struct {
	common      common.Common
	searchField *textfield.Model
	accountPage *account.Model
	listsPage   *lists.Model
	searchPage  *search.Model
	filmPage    *film.Model
	dialog      *dialog.Model
	help        help.Model
	// scrollView  *vlist.Model
}

func New(c common.Common) *Model {

	searchField := textfield.New(c)
	searchField.CharLimit(80)
	searchField.Placeholder("(s)earch for films...")

	m := &Model{
		common:      c,
		searchField: searchField,
		accountPage: account.New(c),
		listsPage:   lists.New(c),
		searchPage:  search.New(c),
		filmPage:    film.New(c),
		dialog:      dialog.New(c, "Quit program?"),
		help:        help.New(),
		// scrollView: vlist.New(c, []tea.Model{
		// 	account.New(c), account.New(c), account.New(c), account.New(c), account.New(c),
		// }),
	}

	m.dialog.Buttons(
		*button.New(c, "Yes", tea.Quit),
		*button.New(c, "No", func() tea.Msg {
			m.dialog.Blur()
			return nil
		}))

	m.SetSize(c.Width, c.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	// wm, hm := ui.getMargins()

	m.common.Width = width
	m.common.Height = height

	// title + " " + searchField = width
	m.searchField.SetSize(width-lipgloss.Width(title)-1, 3)

	contentHeight := height - 3

	m.accountPage.SetSize(util.Max(width/2, 30), contentHeight)

	m.listsPage.SetSize(width, contentHeight)
	m.searchPage.SetSize(width, contentHeight)
	m.filmPage.SetSize(width, contentHeight)

	// m.scrollView.SetSize(width, contentHeight)

	m.help.Width = width
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.AuthState:
		m.common.Global.AuthState.Authed = msg.Authed
		m.common.Global.AuthState.Cookie = msg.Cookie
		m.common.Global.AuthState.User = msg.User
		return m, m.listsPage.Init()
	case tea.WindowSizeMsg:
		frameW, frameH := m.common.Global.Styles.App.GetFrameSize()

		viewW, viewH := msg.Width-frameW, msg.Height-frameH-1 // -1 for help

		m.SetSize(viewW, viewH)

	case tea.KeyMsg:
		if key.Matches(msg, m.common.Global.KeyMap.Quit) {
			if m.dialog.Focused() {
				return m, tea.Quit
			}
			return m, m.dialog.Focus()
		} else if key.Matches(msg, m.common.Global.KeyMap.Back) {
			if m.dialog.Focused() {
				m.dialog.Blur()
			}
		}

		// TODO all other focusables

		// if !m.searchField.Focused() {

		// if key.Matches(msg, m.common.Global.KeyMap.Search) {
		// 	return m, m.searchField.Focus()
		// }

		// }

	}

	var cmd tea.Cmd
	if m.searchField.Focused() {
		_, cmd = m.searchField.Update(msg)
	} else if m.dialog.Focused() {
		_, cmd = m.dialog.Update(msg)
	} else if m.common.Global.AuthState.Authed {
		_, cmd = m.listsPage.Update(msg)
	} else {
		_, cmd = m.accountPage.Update(msg)
	}

	// m.help, cmd = m.help.Update(msg)
	// cmds = append(cmds, cmd)

	return m, cmd
}

func (m *Model) View() string {

	view := strings.Builder{}
	if !m.common.Global.AuthState.Authed {
		// 3 tall to match search bar + fullwidth to allow centering accountPage view
		rightPad := util.Max(m.common.Width-ansi.PrintableRuneWidth(title), 0)
		appBar := "\n" + title + strings.Repeat(" ", rightPad) + "\n"

		centered := lipgloss.JoinVertical(lipgloss.Center, appBar, m.accountPage.View())
		view.WriteString(centered)
	} else {
		appBar := lipgloss.JoinHorizontal(lipgloss.Center, title, " ", m.searchField.View())
		view.WriteString(appBar)
		view.WriteString("\n")

		view.WriteString(m.listsPage.View())
	}

	vGap := m.common.Height - lipgloss.Height(view.String())

	if vGap > 0 {
		view.WriteString(strings.Repeat("\n", vGap))
	}
	// view.WriteString("\n")
	view.WriteString(m.help.View(m.common.Global.KeyMap))

	app := view.String()

	if m.dialog.Focused() {
		dialogView := m.dialog.View()

		dialogW := lipgloss.Width(dialogView)
		dialogH := lipgloss.Height(dialogView)

		xOffset := util.Max((m.common.Width-dialogW)/2, 0)
		yOffset := util.Max((m.common.Height-dialogH)/2-3, 0)

		app = util.RenderOverlay(app, m.dialog.View(), xOffset, yOffset)
	}

	return m.common.Global.Styles.App.Render(app)

}
