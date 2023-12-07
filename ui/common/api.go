package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
)

// const ReviewBase = "http://localhost:8080"
const ReviewBase = "https://review-api.fly.dev"

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

type ByStatusAndUpdate []Review

func (a ByStatusAndUpdate) Len() int      { return len(a) }
func (a ByStatusAndUpdate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByStatusAndUpdate) Less(i, j int) bool {
	if a[i].Status != a[j].Status {
		return a[i].Status > a[j].Status
	}

	if a[i].Updated_at.Equal(a[j].Updated_at) {
		return a[i].Created_at.After(a[j].Updated_at)
	}
	return a[i].Updated_at.After(a[j].Updated_at)
}

type Paged[T Review | Film] struct {
	Results       []T
	Page          int
	Total_Pages   int
	Total_Results int
}

type responseData interface {
	Film | Review | Paged[Film] | Paged[Review]
}

type fetchCallback[T responseData] func(data T, err error) tea.Msg

// Convenience Fetch() wrapper for most common usecase
func Get[T responseData](g Global, url string, callback fetchCallback[T]) tea.Cmd {
	return Fetch[T](g, "GET", url, nil, callback)
}

func Fetch[T responseData](g Global, method string, url string, body map[string]interface{}, callback fetchCallback[T]) tea.Cmd {

	var rawbody []byte
	var err error
	if body != nil {
		rawbody, err = json.Marshal(body)
		if err != nil {
			return nil
		}
	}

	req, err := retryablehttp.NewRequest(method, url, rawbody)
	req.AddCookie(&http.Cookie{Name: "id", Value: g.AuthState.Cookie})
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return nil
	}

	return func() tea.Msg {
		var data T

		resp, err := g.HttpClient.Do(req)

		if err != nil {
			return func() tea.Msg { return callback(data, err) }
		}

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return func() tea.Msg { return callback(data, errors.New(fmt.Sprint(resp.StatusCode))) }
		}

		if resp.StatusCode != 204 {
			err = json.NewDecoder(resp.Body).Decode(&data)
		}

		if err != nil {
			return func() tea.Msg { return callback(data, err) }
		}

		return func() tea.Msg { return callback(data, nil) }
	}
}

const filmEndpoint = "https://api.themoviedb.org/3/movie/"

func GetFilmCmd(g Global, filmId int) tea.Cmd {
	g.FilmCache.SetLoading(filmId)
	url := (filmEndpoint + strconv.Itoa(filmId) + "?api_key=" + g.Config.TMDB_API_KEY)
	return Get[Film](g, url, func(data Film, err error) tea.Msg {
		if err != nil {
			g.FilmCache.Delete(filmId)
		} else {
			g.FilmCache.Set(filmId, data)
		}
		return nil
	})
}

const filmReviewEndpoint = ReviewBase + "/reviews?category=Film"

func GetMyFilmReviewCmd(g Global, filmId int, callback fetchCallback[Paged[Review]]) tea.Cmd {
	url := filmReviewEndpoint +
		"&tmdb_id=" + strconv.Itoa(filmId) +
		"&user_id=" + strconv.Itoa(g.AuthState.User.Id)
	return Get[Paged[Review]](g, url, callback)
}
