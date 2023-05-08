package util

import (
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// unsafe
type Stack[T any] []T

func (s *Stack[T]) Push(v T) {
	*s = append(*s, v)
}

func (s *Stack[T]) Pop() T {
	res := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return res
}

func (s *Stack[T]) Top() T {
	return (*s)[len(*s)-1]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

type AThing struct {
	start int
	end   int
}

// TODO DOES NOT WORK USE WITH CAUTION
func RenderOverlay(parentView, overlayView string, top, left int) string {

	// // return parentView + fmt.Sprintf("\033[%d;%dH%s", top, right, overlayView)
	stack := make(Stack[AThing], 0)
	var escapeStart int
	var escapeEnd int
	var isReset bool
	var n int

	for i, r := range parentView {
		if r == ansi.Marker {
			escapeStart = i
			isReset = true
		} else if escapeStart > escapeEnd {
			// reset sequence = '\xb1', '[', '0', 'm'
			if i == escapeStart+2 {
				if r != '0' {
					isReset = false
				}
			}
			// isTerminal means isAlphabetic, in lipgloss's case 'm'
			if ansi.IsTerminator(r) {
				if i == escapeStart+3 && isReset {
					stack.Pop()
				} else {
					escapeEnd = i + 1
					stack.Push(AThing{
						escapeStart, escapeEnd,
					})
				}
			}
		} else {
			n += runewidth.RuneWidth(r)

			if n >= left {
				// push stack # of resets
			}

			if n >= left+len(strings.Split(overlayView, "\n")[0]) {
				// push stack

				return "something"
			}
		}
	}

	return "something else"
	// ;m
	// var CSI = termenv.CSI
	// this
	// return fmt.Sprintf("%s%sm%s%sm", CSI, seq, s, CSI+ResetSeq)

	// for _, c := range s {
	// 	if c == Marker {
	// 		// ANSI escape sequence
	// 		ansi = true
	// 	} else if ansi {
	// 		if IsTerminator(c) {
	// 			// ANSI sequence terminated
	// 			ansi = false
	// 		}
	// 	} else {
	// 		n += runewidth.RuneWidth(c)
	// 	}
	// }

	// parentLines := strings.Split(parentView, "\n")
	// parentWidth := runewidth.StringWidth(parentLines[0])
	// // parentHeight := len(parentLines)

	// overlayLines := strings.Split(overlayView, "\n")
	// overlayHeight := len(overlayLines)

	// lipgloss.Width(parentView)

	// for j := 0; j < overlayHeight; j++ {
	// 	overlayWidth := lipgloss.Width(overlayLines[j])
	// 	parentLine := parentLines[j+top]

	// 	diff := overlayWidth + right - lipgloss.Width(parentLine)
	// 	if diff > 0 {
	// 		parentLine = fmt.Sprintf("%s%s", parentLine, strings.Repeat(" ", diff))
	// 	}

	// 	var displayWidth int
	// 	var start int
	// 	var end int

	// 	var padLeft int
	// 	var padRight int

	// 	// var ansiStart int
	// 	// var ansiEnd int

	// 	// TODO save last seen ANSI sequence that terminates inside overlay and append after overlay
	// 	// TODO last ANSI sequence inside that does not terminate and move after

	// 	pl := []rune(parentLine)

	// 	var escapeStart int
	// 	var escapeEnd int

	// 	var i int
	// 	for i = 0; i < len(pl); i++ {

	// 		r := pl[i]
	// 		if r == ansi.Marker {
	// 			escapeStart = i
	// 			continue
	// 		} else if escapeStart > escapeEnd {
	// 			if ansi.IsTerminator(r) {
	// 				escapeEnd = i + 1
	// 				if displayWidth >= right {
	// 					break
	// 				}
	// 			}
	// 			continue
	// 		}

	// 		displayWidth += runewidth.RuneWidth(r)

	// 		if displayWidth >= right {
	// 			start = i
	// 			if displayWidth > right {
	// 				padLeft = 1
	// 			}

	// 			if escapeEnd > escapeStart {
	// 				break
	// 			}
	// 		}

	// 	}

	// 	t1 := string(pl[:max(escapeEnd, start)])

	// 	for i = i + 1; i < len(pl); i++ {

	// 		r := pl[i]
	// 		if r == ansi.Marker {
	// 			escapeStart = i
	// 			continue
	// 		} else if escapeStart > escapeEnd {
	// 			if ansi.IsTerminator(r) {
	// 				escapeEnd = i + 1
	// 				if displayWidth >= right+overlayWidth {
	// 					break
	// 				}
	// 			}
	// 			continue
	// 		}

	// 		displayWidth += runewidth.RuneWidth(r)

	// 		if displayWidth >= right+overlayWidth {
	// 			end = i + 1
	// 			if displayWidth > right+overlayWidth {
	// 				padRight = 1
	// 			}

	// 			if escapeEnd > escapeStart {
	// 				break
	// 			}
	// 		}

	// 	}

	// 	t2 := string(pl[min(escapeEnd, end):])
	// 	// for bi, r := range parentLine {

	// 	// 	if r == ansi.Marker {
	// 	// 		isAnsi = true
	// 	// 		if !leftDone {
	// 	// 			ansiStart = bi
	// 	// 		}
	// 	// 	} else if isAnsi {
	// 	// 		if ansi.IsTerminator(r) {
	// 	// 			isAnsi = false
	// 	// 		} else {
	// 	// 			if !leftDone {
	// 	// 				ansiEnd = bi
	// 	// 			}
	// 	// 		}
	// 	// 	} else {

	// 	// 		displayWidth += runewidth.RuneWidth(r)

	// 	// 		if !leftDone && displayWidth >= right {
	// 	// 			start = bi
	// 	// 			if displayWidth > right {
	// 	// 				padLeft = 1
	// 	// 			}
	// 	// 			leftDone = true
	// 	// 		}

	// 	// 		if leftDone && displayWidth >= right+overlayWidth {
	// 	// 			end = bi + utf8.RuneLen(r)
	// 	// 			if displayWidth > right+overlayWidth {
	// 	// 				padRight = 1
	// 	// 			}
	// 	// 			break
	// 	// 		}
	// 	// 	}
	// 	// }

	// 	// if ansiEnd < ansiStart {
	// 	// parentLine[ansiStart]
	// 	// }

	// 	parentLines[j+top] = fmt.Sprintf("%s%*s%*s", (t1), padLeft, overlayLines[j], padRight, (t2))
	// }

	// return strings.Join(parentLines, "\n")
}
