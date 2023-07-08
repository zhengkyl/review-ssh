package reviewlist

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	top  = "▀"
	bot  = "▄"
	full = "█"
)

var (
	scrollBorder = lipgloss.RoundedBorder()
	scrollStyle  = lipgloss.NewStyle().Border(scrollBorder, true)
)

func renderScrollbar(height, positions, pos int) string {
	// sos := pos
	innerHeight := height - 2 // top/bot border
	thumbPositions := innerHeight * 2

	if positions <= 1 {
		// todo
	}

	thumbHeight := 1 // half char units

	if positions < thumbPositions {
		// add halfblocks equal to diff
		thumbHeight += thumbPositions - positions
	} else if positions > thumbPositions {
		// convert pos to be a ratio of thumbPositions - 1
		// Only max when pos + 1 == positions (last item visible)
		pos = (pos + 1) * (thumbPositions - 1) / positions
	}

	endPos := pos + thumbHeight - 1

	thumbStartIndex := pos / 2
	thumbEndIndex := endPos / 2

	sb := strings.Builder{}
	for i := 0; i < innerHeight; i++ {

		if i == thumbStartIndex {

			if pos%2 == 1 {
				sb.WriteString(bot)
			} else if thumbHeight == 1 {
				sb.WriteString(top)
			} else {
				sb.WriteString(full)
			}

		} else if i == thumbEndIndex {
			if endPos%2 == 0 {
				sb.WriteString(top)
			} else {
				sb.WriteString(full)
			}
		} else if i > thumbStartIndex && i < thumbEndIndex {
			sb.WriteString(full)
		} else {
			sb.WriteString(" ")
		}

		if i != innerHeight-1 {
			sb.WriteString("\n")
		}
	}

	return scrollStyle.Render(sb.String())
}
