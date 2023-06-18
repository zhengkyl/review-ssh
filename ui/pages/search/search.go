package search

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type Model struct {
	props common.Props
	input textinput.Model
	list  list.Model
}

// const film_url = "https://review-api.fly.dev/search/Film"
// const show_url = "https://review-api.fly.dev/search/Show"

type itemJson struct {
	Id           int
	Title        string
	Overview     string
	Poster_path  string
	Release_date string
}

func New(props common.Props) *Model {

	input := textinput.New()
	input.Placeholder = "Search for films and shows..."
	input.Focus()
	input.CharLimit = 80

	m := &Model{
		input: input,
		props: props,
		list:  list.New([]list.Item{}, itemDelegate{}, 0, 0),
	}

	m.SetSize(props.Width, props.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	m.list.SetSize(width, height)
	// wm, hm := m.getMargins()

}

func (m *Model) getMargins() (wm, hm int) {
	wm = 0
	hm = 0

	return
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds = append(cmds, getSearchCmd(m.props.Global.HttpClient, m.input.Value()))
		}

	case []list.Item:
		cmds = append(cmds, m.list.SetItems(msg))
		for _, i := range msg {
			var j = i.(item)
			cmds = append(cmds, j.poster.Init(), j.buttons.Init())
		}

	case error:
		return m, nil
	}

	var cmd tea.Cmd
	// Necessary b/c bubbles component
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	// Necessary b/c bubbles component
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	var view string

	wm, _ := m.getMargins()

	ss := lipgloss.NewStyle().Width(m.props.Width - wm)
	view = ss.Render(m.input.View())
	view += ss.Render(m.list.View())

	return view
}
