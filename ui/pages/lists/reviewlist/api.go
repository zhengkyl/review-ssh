package reviewlist

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
)

const filmEndpoint = "https://api.themoviedb.org/3/movie/"
const showEndpoint = "https://api.themoviedb.org/3/tv/"

func getFilmCmd(g common.Global, filmId int) tea.Cmd {
	url := (filmEndpoint + strconv.Itoa(filmId) + "?api_key=" + g.Config.TMDB_API_KEY)
	return common.GetCmd[common.Film](g.HttpClient, url)
}

func getShowCmd(g common.Global, showId int) tea.Cmd {
	url := (showEndpoint + strconv.Itoa(showId) + "?api_key=" + g.Config.TMDB_API_KEY)
	return common.GetCmd[common.Show](g.HttpClient, url)
}
