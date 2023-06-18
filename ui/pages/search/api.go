package search

import (
	"encoding/json"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/components/poster"
)

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
				r.Release_date,
				poster.New(common.Props{Width: POSTER_WIDTH, Height: POSTER_HEIGHT}, "https://image.tmdb.org/t/p/w200"+r.Poster_path),
				NewButtons(common.Props{Width: 0, Height: 0}),
			}
			itemResults = append(itemResults, i)
		}
		return itemResults
	}
}
