package common

import (
	"encoding/json"
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

type Show struct {
	Id             int
	Name           string
	Overview       string
	Poster_path    string
	First_air_date string
}

type Review struct {
	User_id  int            //`json:"user_id"`
	Tmdb_id  int            //`json:"tmdb_id"`
	Category enums.Category //`json:"category"`
	Season   int            //`json:"season"` only for shows
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
	// Season     int     //`json:"season"`
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

// tmdb_id overlaps for movies and tv series and is not unique per season
func (r Review) Key() int {
	switch r.Category {
	case enums.Film:
		return -r.Tmdb_id
	case enums.Show:
		return r.Tmdb_id*100 + r.Season
	default:
		// TODO error?
		return 0
	}
}

type Resource interface {
	Film | Show
}
type GetResponse[T Resource] struct {
	Ok   bool
	Data T
	Err  string
}

func GetCmd[T Resource](client *retryablehttp.Client, url string) tea.Cmd {
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
