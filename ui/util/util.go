package util

import (
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/muesli/ansi"
)

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
		// Adding previous nonterminating control sequences
		sb := strings.Builder{}

		growth := stackBytes + len(line)
		sb.Grow(growth)

		for _, seq := range stack {
			sb.WriteString(lines[seq.line][seq.start:seq.end])
		}
		sb.WriteString(line)

		shouldClearPrefix := stackBytes > 0

		// Detecting nonterminiated control sequences in line
		seqStart := -1
		seqEnd := 0

		var width int
		var r rune
		for j := 0; j < len(line); j += width {
			r, width = utf8.DecodeRuneInString(line[j:])

			if r == ansi.Marker {
				seqStart = j
			} else if seqStart >= seqEnd {
				if !ansi.IsTerminator(r) {
					continue
				}

				seqEnd = j + 1

				if j == seqStart+3 &&
					// Control sequences are ASCII so, string index is fine
					line[seqStart+1] == '[' &&
					line[seqStart+2] == '0' &&
					line[seqStart+3] == 'm' {

					stack.Clear()
					stackBytes = 0

					shouldClearPrefix = false
				} else {
					stack.Push(LineSeqBounds{i, seqStart, seqEnd})
					stackBytes += seqEnd - seqStart
				}
			}
		}

		if shouldClearPrefix {
			sb.WriteString(resetSeq)
		}
		styledLines[i] = sb.String()
	}

	return styledLines
}

func RenderOverlay(parentView, overlayView string, left, top int) string {
	parentLines := strings.Split(parentView, "\n")
	overlayLines := styleByLine(overlayView)

	finalLines := make([]string, len(parentLines))
	stack := make(Stack[LineSeqBounds], 0)

	for i, parentLine := range parentLines {

		n := 0

		// Detect when seqStart > seqEnd, even when 0 is valid
		seqStart := -1
		seqEnd := 0

		overlayStart := false
		overlayEnd := false

		overlayIndex := i - top
		shouldOverlay := overlayIndex >= 0 && overlayIndex < len(overlayLines)

		// lineWidth := ansi.PrintableRuneWidth(parentLine)
		// if shouldOverlay && lineWidth < left {
		// 	parentLine += strings.Repeat(" ", left-lineWidth)
		// }

		sb := strings.Builder{}

		var width int
		var r rune
		for j := 0; j < len(parentLine); j += width {
			r, width = utf8.DecodeRuneInString(parentLine[j:])

			if r == ansi.Marker {
				seqStart = j
			} else if seqStart >= seqEnd {
				if !ansi.IsTerminator(r) {
					continue
				}

				seqEnd = j + 1

				if j == seqStart+3 &&
					parentLine[seqStart+1] == '[' &&
					parentLine[seqStart+2] == '0' &&
					parentLine[seqStart+3] == 'm' {

					stack.Clear()
				} else {
					stack.Push(LineSeqBounds{i, seqStart, seqEnd})
				}

			} else if shouldOverlay {
				n += runewidth.RuneWidth(r)

				if !overlayStart {
					if n < left {
						continue
					}
					overlayStart = true

					if left != 0 {
						if n == left {
							sb.WriteString(parentLine[:j+width])
						} else if n > left {
							sb.WriteString(parentLine[:j] + " ")
						}

						sb.WriteString(resetSeq)
					}

					sb.WriteString(overlayLines[overlayIndex])
				} else if !overlayEnd {
					right := left + ansi.PrintableRuneWidth(overlayLines[overlayIndex])

					if n < right {
						continue
					}
					overlayEnd = true

					for _, seq := range stack {
						sb.WriteString(parentLines[seq.line][seq.start:seq.end])
					}

					if n == right {
						sb.WriteString(parentLine[j+width:])
					} else if n > right {
						sb.WriteString(" " + parentLine[j+width:])
					}
				}
			}

		}

		if shouldOverlay {
			if !overlayStart {
				sb.WriteString(parentLine)
				sb.WriteString(strings.Repeat(" ", left-n))
				sb.WriteString(overlayLines[overlayIndex])
			}
			finalLines[i] = sb.String()
		} else {
			finalLines[i] = parentLine
		}

	}

	return strings.Join(finalLines, "\n")
	// reference
	// var CSI = termenv.CSI
	// return fmt.Sprintf("%s%sm%s%sm", CSI, seq, s, CSI+ResetSeq)
}
