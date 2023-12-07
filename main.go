package main

import (
	"log"
	"os"

	"github.com/zhengkyl/review-ssh/server"
)

func main() {
	tmdbKey, ok := os.LookupEnv("TMDB_API_KEY")
	if !ok {
		log.Fatal("TMDB_API_KEY missing")
	}

	server.RunServer(tmdbKey)
}

// func runLocal() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	tmdbKey, ok := os.LookupEnv("TMDB_API_KEY")
// 	if !ok {
// 		log.Fatal("TMDB_API_KEY missing")
// 	}

// 	httpClient := retryablehttp.NewClient()
// 	httpClient.Logger = nil

// 	c := common.Props{
// 		Global: common.Global{
// 			AuthState: &common.AuthState{
// 				Authed: false,
// 			},
// 			Config: common.Config{
// 				TMDB_API_KEY: tmdbKey,
// 			},

// 			ReviewMap:  map[int]common.Review{},
// 			FilmCache:  common.Cache[common.Film]{},
// 			KeyMap:     keymap.DefaultKeyMap(),
// 			HttpClient: httpClient,
// 		},
// 	}

// 	p := tea.NewProgram(ui.New(c), tea.WithAltScreen())

// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("L + R, Kyle fix your code: %v", err)
// 		os.Exit(1)
// 	}
// }
