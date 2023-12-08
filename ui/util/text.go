package util

import (
	"github.com/mattn/go-runewidth"
)

func TruncAndPadUnicode(text string, length int) string {
	if length < 1 {
		return ""
	}

	len := runewidth.StringWidth(text)

	if len == length {
		return text
	}

	if length > len {
		return runewidth.FillRight(text, length)
	}

	var trunc string
	if text[length-2] == ' ' {
		trunc = runewidth.Truncate(text, length, "… ")
	} else {
		trunc = runewidth.Truncate(text, length, "…")
	}
	return runewidth.FillRight(trunc, length)
}
