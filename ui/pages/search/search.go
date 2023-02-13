package search

import (
	"encoding/json"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
)

type SearchModel struct {
	common     common.Common
	httpClient *retryablehttp.Client
	input      textinput.Model
	list       list.Model
}

// const film_url = "https://review-api.fly.dev/search/Film"
// const show_url = "https://review-api.fly.dev/search/Show"

type itemJson struct {
	Id           int
	Title        string
	Overview     string
	Poster_Path  string
	Release_Date string
}
type searchResponse struct {
	Results []itemJson
	// unused fields
	// Page          int
	// Total_Pages   int
	// Total_Results int
}

func getSearchCmd(client *retryablehttp.Client, query string) tea.Cmd {

	return func() tea.Msg {
		// resp, err := client.Get(fmt.Sprintf("%s?query=%s", film_url, query))
		resp, err := client.Get("https://review-api.fly.dev/search/Film?query=" + query)
		if err != nil {
			return []item{}
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return []item{}
		}

		var searchResponse searchResponse
		err = json.Unmarshal(body, &searchResponse)
		if err != nil {
			return []item{}
		}

		var itemResults []list.Item

		for _, r := range searchResponse.Results {
			i := item{
				r.Id,
				r.Title,
				r.Overview,
				r.Release_Date,
				poster.New(common.Common{Width: POSTER_WIDTH, Height: POSTER_HEIGHT}, "https://image.tmdb.org/t/p/w200"+r.Poster_Path),
				NewButtons(common.Common{Width: 0, Height: 0}),
			}
			itemResults = append(itemResults, i)
		}
		return itemResults
	}
}

func New(common common.Common, httpClient *retryablehttp.Client) *SearchModel {

	input := textinput.New()
	input.Placeholder = "Search for movies and shows..."
	input.Focus()
	input.CharLimit = 80

	m := &SearchModel{
		input:      input,
		common:     common,
		list:       list.New([]list.Item{}, itemDelegate{}, 0, 0),
		httpClient: httpClient,
	}

	m.SetSize(common.Width, common.Height)

	return m
}

func (m *SearchModel) SetSize(width, height int) {
	m.common.Width = width
	m.common.Height = height

	m.list.SetSize(width, height)
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds = append(cmds, getSearchCmd(m.httpClient, m.input.Value()))
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

	return view
}
