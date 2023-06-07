package styles

import "github.com/charmbracelet/lipgloss"

// This defines a struct shared between ALL components
// Make sure to use .Copy() before modifying any style for a component
type Styles struct {
	App         lipgloss.Style
	SearchInput lipgloss.Style
}

func DefaultStyles() *Styles {

	return &Styles{
		App:         lipgloss.NewStyle().Margin(0, 1),
		SearchInput: lipgloss.NewStyle(),
	}
}
