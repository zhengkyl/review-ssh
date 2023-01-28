package search

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/image"
)

type SearchModel struct {
	common  common.Common
	input   textinput.Model
	results []item
	list    list.Model
	t1      bool
	t2      bool
}

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

const film_url = "https://review-api.fly.dev/search/Film"
const show_url = "https://review-api.fly.dev/search/Show"

type item struct {
	id           int
	title        string
	overview     string
	release_date string
	image        *image.ImageModel
}

type itemJson struct {
	Id          int
	Title       string
	Overview    string
	Poster_Path string
	ReleaseDate string
}

// func (i item) Title() string       { return i.title }
// func (i item) Description() string { return i.overview }
func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int  { return 20 }
func (d itemDelegate) Spacing() int { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, listItem := range m.Items() {
		i, ok := listItem.(item)
		if !ok {
			return nil
		}

		// var cmd tea.Cmd

		_, cmd := i.image.Update(msg)

		// i.image = imageM.(image.ImageModel)
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// str := fmt.Sprintf("%d. %s", index+1, i.title)
	str := lipgloss.JoinHorizontal(0, i.image.View(), i.title)

	fn := itemStyle.Render
	// if index == m.Index() {
	// 	fn = func(s string) string {
	// 		return selectedItemStyle.Render("> " + s)
	// 	}
	// }

	fmt.Fprint(w, fn(str))
}

type searchMsg []item

type searchResponse struct {
	Results []itemJson
}

func search(query string) tea.Cmd {

	return func() tea.Msg {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		// resp, err := client.Get(fmt.Sprintf("%s?query=%s", film_url, query))
		resp, err := client.Get("https://review-api.fly.dev/search/Film?query=" + query)

		if err != nil {
			return searchMsg([]item{{id: 1, title: "here" + err.Error()}})
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return searchMsg([]item{{id: 2, title: "there" + err.Error()}})
		}

		var searchResponse searchResponse
		// json.NewDecoder(resp.Body).Decode(&film)
		err = json.Unmarshal(body, &searchResponse)

		if err != nil {
			return searchMsg([]item{{id: 3, title: "everywhere " + err.Error()}})
		}

		var itemResults []item

		for _, x := range searchResponse.Results {
			itemResults = append(itemResults, item{
				x.Id,
				x.Title,
				x.Overview,
				x.ReleaseDate,
				image.New(common.Common{}, "https://image.tmdb.org/t/p/w200"+x.Poster_Path),
			})
		}
		return searchMsg(itemResults)
	}
}

func New(common common.Common) *SearchModel {

	input := textinput.New()
	input.Placeholder = "Search for movies and shows..."
	input.Focus()
	input.CharLimit = 80

	m := &SearchModel{
		input:  input,
		common: common,
		list:   list.New([]list.Item{}, itemDelegate{}, 50, 20),
	}

	m.SetSize(common.Width, common.Height)

	return m
}

func (m *SearchModel) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	// wm, hm := m.getMargins()

}

func (m *SearchModel) getMargins() (wm, hm int) {
	wm = 0
	hm = 0

	return
}

func (m *SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.t2 = true
			cmds = append(cmds, search(m.input.Value()))
		}
	case searchMsg:
		m.t1 = true
		m.results = msg
		b := make([]list.Item, len(msg))
		for i := range msg {
			b[i] = msg[i]
		}
		cmd := m.list.SetItems(b)
		return m, cmd
		// cmds = append(cmds, cmd)

	case error:
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *SearchModel) View() string {
	var view string

	wm, _ := m.getMargins()

	ss := lipgloss.NewStyle().Width(m.common.Width - wm)
	view = ss.Render(m.input.View())
	view += ss.Render(m.list.View())

	if m.t1 {
		view += "t1 true"
	}
	if m.t2 {
		view += "t2 true"
	}

	view += fmt.Sprintf("here are my %v", m.results)

	return view
}

// var cmd tea.Cmd

// m.searchInput, cmd = m.searchInput.Update(msg)
