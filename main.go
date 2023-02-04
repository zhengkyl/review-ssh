package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui"
)

func main() {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	p := tea.NewProgram(ui.New(httpClient), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("L + R, Kyle fix your code: %v", err)
		os.Exit(1)
	}
}
