package common

import "github.com/charmbracelet/lipgloss"

var (
	ratings = []string{
		"░░ ░░ ░░",
		"██ ░░ ░░",
		"░░ ██ ░░",
		"██ ██ ░░",
		"░░ ░░ ██",
		"██ ░░ ██",
		"░░ ██ ██",
		"██ ██ ██",
	}
	t  = "██\n \n░░"
	f  = "░░\n \n██"
	tt = "-\n \n "
	ff = " \n \n-"
	tf = " \n\\\n "
	ft = " \n/\n "

	activeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func RenderThinRating(before, during, after bool) string {
	ratingIndex := 0
	if before {
		ratingIndex += 1
	}
	if during {
		ratingIndex += 2
	}
	if after {
		ratingIndex += 4
	}
	return ratings[ratingIndex]
}

func RenderThickRating(before, during, after bool) string {
	var sections []string

	if before {
		sections = append(sections, activeStyle.Render(t))
	} else {
		sections = append(sections, activeStyle.Render(f))
	}

	sections = append(sections, tt)

	if during {
		sections = append(sections, t)
	} else {
		sections = append(sections, f)
	}

	sections = append(sections, tf)

	if after {
		sections = append(sections, t)
	} else {
		sections = append(sections, f)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, sections...)
}
