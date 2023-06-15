package common

import "time"

type User struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type Movie struct {
	Id           int
	Title        string
	Overview     string
	Poster_path  string
	Release_date string
}

// type Show struct {
// 	Id             int
// 	Name           string
// 	Overview       string
// 	Poster_path    string
// 	First_air_date string
// }

type Review struct {
	User_id    int       //`json:"user_id"`
	Tmdb_id    int       //`json:"tmdb_id"`
	Category   string    //`json:"category"`
	Status     string    //`json:"status"`
	Text       string    //`json:"text"`
	Fun_before bool      //`json:"fun_before"`
	Fun_during bool      //`json:"fun_during"`
	Fun_after  bool      //`json:"fun_after"`
	Created_at time.Time //`json:"created_at"`
	Updated_at time.Time //`json:"updated_at"`
	// Season     int     //`json:"season"`
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

func (a ByUpdatedAt) Len() int           { return len(a) }
func (a ByUpdatedAt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByUpdatedAt) Less(i, j int) bool { return a[i].Updated_at.Before(a[j].Updated_at) }
