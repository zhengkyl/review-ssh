package reviewlist

import (
	"encoding/json"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

const endpoint = "https://api.themoviedb.org/3/movie/"

type res struct {
	ok  bool
	err string
}

func getMovie(client *retryablehttp.Client, apiKey string, movieId int) tea.Msg {
	resp, err := client.Get(endpoint + strconv.Itoa(movieId) + "?api_key=" + apiKey)

	if err != nil {
		return res{false, err.Error()}
	}

	if resp.StatusCode != 200 {
		return res{false, "Something went wrong."}
	}

	var response common.Movie
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return res{false, err.Error()}
	}

	return response
}
