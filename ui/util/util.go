package util

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
)

// TODO DOES NOT WORK USE WITH CAUTION
func RenderOverlay(parentView, overlayView string, top, right int) string {

	parentLines := strings.Split(parentView, "\n")
	// parentWidth := runewidth.StringWidth(parentLines[0])
	// parentHeight := len(parentLines)

	overlayLines := strings.Split(overlayView, "\n")
	overlayWidth := runewidth.StringWidth(overlayLines[0])
	overlayHeight := len(overlayLines)

	for i := 0; i < overlayHeight; i++ {
		line := parentLines[i+top]

		var displayWidth int
		var start int
		var end int

		var padLeft int
		var padRight int

		leftDone := false

		for bi, r := range line {
			displayWidth += runewidth.RuneWidth(r)

			if !leftDone && displayWidth >= right {
				start = bi
				if displayWidth > right {
					padLeft = 1
				}
				leftDone = true
			}

			if leftDone && displayWidth >= right+overlayWidth {
				end = bi + utf8.RuneLen(r)
				if displayWidth > right+overlayWidth {
					padRight = 1
				}
				break
			}
		}

		parentLines[i+top] = fmt.Sprintf("%s%*s%*s", line[:start], padLeft, overlayLines[i], padRight, line[end:])
	}

	return strings.Join(parentLines, "\n")
}
