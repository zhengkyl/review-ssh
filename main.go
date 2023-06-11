package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/keymap"
	"github.com/zhengkyl/review-ssh/ui/styles"
)

func main() {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	tmdbKey, ok := os.LookupEnv("TMDB_API_KEY")
	if !ok {
		panic("TMDB_API_KEY missing")
	}

	c := common.Common{
		Global: common.Global{
			Config: common.Config{
				TMDB_API_KEY: tmdbKey,
			},
			Styles:     styles.DefaultStyles(),
			KeyMap:     keymap.DefaultKeyMap(),
			HttpClient: httpClient,
		},
	}

	p := tea.NewProgram(ui.New(c), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("L + R, Kyle fix your code: %v", err)
		os.Exit(1)
	}
}
