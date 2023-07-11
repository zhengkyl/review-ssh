package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/textfield"
	"github.com/zhengkyl/review-ssh/ui/components/vlist"
	"github.com/zhengkyl/review-ssh/ui/pages/search/filmitem"
)

type Model struct {
	props       common.Props
	searchField *textfield.Model
	list        *vlist.Model
}

func New(p common.Props) *Model {

	searchField := textfield.New(p)
	searchField.CharLimit(80)
	searchField.Placeholder("(s)earch for films...")

	m := &Model{
		props:       p,
		searchField: searchField,
		list:        vlist.New(p, 4),
	}

	m.SetSize(p.Width, p.Height)

	return m
}

func (m *Model) SetSize(width, height int) {
	m.props.Width = width
	m.props.Height = height

	// title + " " + searchField = width
	m.searchField.SetSize(width, 3)

	m.list.SetSize(width, height)
	// wm, hm := m.getMargins()

}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (common.Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd := common.Get[common.Paged[common.Film]](m.props.Global.HttpClient, endpoint+"?query="+m.searchField.Value(), func(data common.Paged[common.Film], err error) tea.Msg {
				if err == nil {
					items := make([]common.Focusable, 0, len(data.Results))
					for _, film := range data.Results {
						item := filmitem.New(m.props, film)
						cmds = append(cmds, item.Init())
					}
					m.list.SetItems(items)

				}
				return nil
			})
			cmds = append(cmds, cmd)
		}
	}

	var cmd tea.Cmd

	_, cmd = m.searchField.Update(msg)
	cmds = append(cmds, cmd)

	_, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) ViewSearchbar() string {
	return m.searchField.View()
}

func (m *Model) View() string {
	sb := strings.Builder{}

	// ss := lipgloss.NewStyle().Width(m.props.Width - wm)
	// view = ss.Render(m.input.View())
	sb.WriteString(m.list.View())

	return sb.String()
}
