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

// Pop() is unchecked
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

func (s *Stack[T]) Clear() {
	*s = (*s)[:0]
}

func (s *Stack[T]) IsEmpty() bool {
	return len(*s) == 0
}

type LineSeqBounds struct {
	line  int
	start int
	end   int
}

const resetSeq = "\033[0m"

func styleByLine(view string) []string {
	lines := strings.Split(view, "\n")

	if len(lines) == 1 {
		return lines
	}

	styledLines := make([]string, len(lines))

	stack := make(Stack[LineSeqBounds], 0)
	stackBytes := 0

	for i, line := range lines {
		sb := strings.Builder{}

		growth := stackBytes + len(line)
		if stackBytes > 0 {
			growth += len(resetSeq)
		}
		sb.Grow(growth)

		for _, seq := range stack {
			sb.WriteString(lines[seq.line][seq.start:seq.end])
		}
		sb.WriteString(line)

		if stackBytes > 0 {
			sb.WriteString(resetSeq)
		}

		styledLines[i] = sb.String()

		lineRunes := []rune(line)

		var seqStart int
		var seqEnd int

		for j := 0; j < len(lineRunes); j++ {
			r := lineRunes[j]

			if r == ansi.Marker {
				seqStart = j
			} else if seqStart > seqEnd {
				if !ansi.IsTerminator(r) {
					continue
				}

				if j == seqStart+3 &&
					lineRunes[seqStart+1] == '[' &&
					lineRunes[seqStart+2] == '0' &&
					lineRunes[seqStart+3] == 'm' {

					stack.Clear()
					stackBytes = 0
				} else {
					seqEnd = j + 1
					stack.Push(LineSeqBounds{i, seqStart, seqEnd})
					stackBytes += seqEnd - seqStart
				}
			}
		}

	}

	return styledLines
}

// type SeqBounds struct {
// 	start int
// 	end   int
// }

// TODO DOES NOT WORK USE WITH CAUTION
func RenderOverlay(parentView, overlayView string, top, left int) string {

	parentLines := strings.Split(parentView, "\n")

	// TODO benchmark alternatives if []rune() is slow
	// strings.NewReader().ReadRune()
	// utf8.DecodeRuneInString()

	overlayLines := styleByLine(overlayView)

	layeredLines := make([]string, len(parentLines))

	stack := make(Stack[LineSeqBounds], 0)

	for i, parentLine := range parentLines {

		n := 0
		parentRunes := []rune(parentLine)

		var seqStart int
		var seqEnd int

		for j := 0; j < len(parentRunes); j++ {
			r := parentRunes[j]

			if r == ansi.Marker {
				seqStart = j
			} else if seqStart > seqEnd {
				if !ansi.IsTerminator(r) {
					continue
				}

				if j == seqStart+3 &&
					parentRunes[seqStart+1] == '[' &&
					parentRunes[seqStart+2] == '0' &&
					parentRunes[seqStart+3] == 'm' {

					stack.Clear()
				} else {
					seqEnd = j + 1
					stack.Push(LineSeqBounds{i, seqStart, seqEnd})
				}

			} else if i >= top && i < top+len(overlayLines) {
				n += runewidth.RuneWidth(r)

				// TODO insert control sequences

				if n < left {
					continue
				}

				if n > left {

					// parentLines[i][:j+1]
					// resetSeq
					// overlayLines[i]
				}

				if n >= left+len(overlayLines[i]) {
					// parentLines[i][j:]
				}

			}
		}
	}

	return strings.Join(layeredLines, "\n")
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
}
