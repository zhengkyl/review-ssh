package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type SearchModel struct {
	common  common.Common
	input   textinput.Model
	results []item
	list    list.Model
	t1      bool
	t2      bool
}

const film_url = "https://review-api.fly.dev/search/Film"
const show_url = "https://review-api.fly.dev/search/Show"

type item struct {
	title string
	Id    int
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return "place holder desc" }
func (i item) FilterValue() string { return i.title }

type searchMsg []item

type searchResponse struct {
	Results []item
}

func search(query string) tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		// resp, err := client.Get(fmt.Sprintf("%s?query=%s", film_url, query))
		resp, err := client.Get("https://review-api.fly.dev/search/Film?query=" + query)

		if err != nil {
			return searchMsg([]item{{"here" + err.Error(), 1}})
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return searchMsg([]item{{"there" + err.Error(), 2}})
		}

		var searchResponse searchResponse
		// json.NewDecoder(resp.Body).Decode(&film)
		err = json.Unmarshal(body, &searchResponse)

		if err != nil {
			return searchMsg([]item{{"everywhere " + err.Error(), 3}})
		}

		return searchMsg(searchResponse.Results)
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
		list:   list.New([]list.Item{}, list.NewDefaultDelegate(), 50, 20),
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
