package common

import (
	"encoding/json"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
)

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Film struct {
	Id           int
	Title        string
	Overview     string
	Poster_path  string
	Release_date string
}

type Review struct {
	User_id  int            //`json:"user_id"`
	Tmdb_id  int            //`json:"tmdb_id"`
	Category enums.Category //`json:"category"`
	//
	Status     enums.Status //`json:"status"`
	Text       string       //`json:"text"`
	Fun_before bool         //`json:"fun_before"`
	Fun_during bool         //`json:"fun_during"`
	Fun_after  bool         //`json:"fun_after"`
	Created_at time.Time    //`json:"created_at"`
	Updated_at time.Time    //`json:"updated_at"`
}

type ReviewUpdate struct {
	Status     string `json:"status"`
	Text       string `json:"text"`
	Fun_before bool   `json:"fun_before"`
	Fun_during bool   `json:"fun_during"`
	Fun_after  bool   `json:"fun_after"`
}

type ReviewNew struct {
	Tmdb_id  int    `json:"tmdb_id"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

type ByUpdatedAt []Review

func (a ByUpdatedAt) Len() int      { return len(a) }
func (a ByUpdatedAt) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUpdatedAt) Less(i, j int) bool {
	if a[i].Updated_at.Equal(a[j].Updated_at) {
		return a[i].Created_at.After(a[j].Updated_at)
	}
	return a[i].Updated_at.After(a[j].Updated_at)
}

type Paginated interface {
	Review
}

type PageResult[T Paginated] struct {
	Results       []T
	Page          int
	Total_Pages   int
	Total_Results int
}

type Gettable interface {
	PageResult[Review] | Film
}

type GetResponse[T Gettable] struct {
	Ok   bool
	Data T
	Err  string
}

func GetCmd[T Gettable](client *retryablehttp.Client, url string) tea.Cmd {
	return func() tea.Msg {
		resp, err := client.Get(url)

		var data T
		if err != nil {
			return GetResponse[T]{false, data, err.Error()}
		}

		if resp.StatusCode != 200 {
			return GetResponse[T]{false, data, "Something went wrong."}
		}

		err = json.NewDecoder(resp.Body).Decode(&data)

		if err != nil {
			return GetResponse[T]{false, data, err.Error()}
		}

		return GetResponse[T]{true, data, ""}
	}
}

const filmEndpoint = "https://api.themoviedb.org/3/movie/"

func GetFilmCmd(g Global, filmId int) tea.Cmd {
	url := (filmEndpoint + strconv.Itoa(filmId) + "?api_key=" + g.Config.TMDB_API_KEY)
	return GetCmd[Film](g.HttpClient, url)
}

const filmReviewEndpoint = "https://review-api.fly.dev/reviews?category=Film"

func GetMyFilmReviewCmd(g Global, filmId int) tea.Cmd {
	url := filmReviewEndpoint +
		"&tmdb_id=" + strconv.Itoa(filmId) +
		"&user_id=" + strconv.Itoa(g.AuthState.User.Id)
	return GetCmd[PageResult[Review]](g.HttpClient, url)
}
