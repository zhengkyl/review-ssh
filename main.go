package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// var style = lipgloss.NewStyle().
// 	Bold(true).
// 	Foreground(lipgloss.Color("#FAFAFA")).
// 	Background(lipgloss.Color("#7D56F4")).
// 	PaddingTop(2).
// 	PaddingLeft(4).
// 	Width(22)

var (
	testStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(2).
			PaddingLeft(4).
			Width(22)
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	docStyle       = lipgloss.NewStyle().Background(lipgloss.Color("#7D56F4")).Padding(1, 2)
	windowStyle    = lipgloss.NewStyle().BorderForeground(highlightColor)
)

type tab struct {
	name string
	view string
}

type model struct {
	searchInput textinput.Model
	tabs        []tab
	activeTab   int
	// selected    map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	searchInput := textinput.New()
	searchInput.Placeholder = "Search for movies and shows..."
	searchInput.Focus()

	tabs := []tab{
		{
			name: "Search",
		},
		{
			name: "My List",
		},
		{
			name: "Account",
		},
	}

	return model{
		searchInput: searchInput,
		tabs:        tabs,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

			// The "up" and "k" keys move the cursor up
		}
	}

	var cmd tea.Cmd

	m.searchInput, cmd = m.searchInput.Update(msg)
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func (m model) View() string {

	view := strings.Builder{}

	view.WriteString(m.searchInput.View())

	// for i, tab := range m.tabs {

	// }

	// view.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.tabs[m.activeTab]))
	// The footer
	view.WriteString("\nPress q to quit.\n")

	// Send the UI for rendering
	return docStyle.Render(view.String())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("L + R, Kyle fix your code: %v", err)
		os.Exit(1)
	}
}
