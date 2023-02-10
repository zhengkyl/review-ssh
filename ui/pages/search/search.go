package search

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/image"
	"golang.org/x/exp/slices"
)

type SearchModel struct {
	common     common.Common
	httpClient *retryablehttp.Client
	input      textinput.Model
	list       list.Model
}

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().Padding(0, 4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(4).Foreground(lipgloss.Color("170"))
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	// helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

const film_url = "https://review-api.fly.dev/search/Film"
const show_url = "https://review-api.fly.dev/search/Show"

// NOTE: Fullwidth spaces are 2 wide
const POSTER_WIDTH = 4 * 2
const POSTER_HEIGHT = 6

type item struct {
	id           int
	title        string
	overview     string
	release_date string
	image        *image.ImageModel
}

// implement list.Item
func (i item) FilterValue() string { return i.title }

type itemJson struct {
	Id          int
	Title       string
	Overview    string
	Poster_Path string
	ReleaseDate string
}

type itemDelegate struct{}

// implement list.ItemDelegate
func (d itemDelegate) Height() int { return POSTER_HEIGHT }

// implement list.ItemDelegate
func (d itemDelegate) Spacing() int { return 0 }

// implement list.ItemDelegate
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	var cmds []tea.Cmd

	for _, listItem := range m.Items() {
		i := listItem.(item)

		_, cmd := i.image.Update(msg)

		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
}

var textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
var titleStyle = lipgloss.NewStyle().Bold(true)
var contentStyle = lipgloss.NewStyle().MarginLeft(2)

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

// implement list.ItemDelegate
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(item)

	contentWidth := m.Width() - itemStyle.GetHorizontalFrameSize() - POSTER_WIDTH - contentStyle.GetHorizontalFrameSize() - 10

	// Subtract 15 to account for long word causing early newline.
	desc := ellipsisText(i.overview, contentWidth*2-15)

	str := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render(i.title), textStyle.Width(contentWidth).Render(desc))

	str = contentStyle.Render(str)

	str = lipgloss.JoinHorizontal(lipgloss.Top, i.image.View(), str)

	if index == m.Index() {
		str = lipgloss.JoinHorizontal(lipgloss.Left, "> ", str)
		str = selectedItemStyle.Render(str)
	} else {
		str = itemStyle.Render(str)
	}

	fmt.Fprint(w, str)
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
				r.ReleaseDate,
				image.New(common.Common{Width: POSTER_WIDTH, Height: POSTER_HEIGHT}, "https://image.tmdb.org/t/p/w200"+r.Poster_Path),
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
			cmds = append(cmds, j.image.Init())
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
