package util

import "strings"

func TruncOrPadASCII(text string, length int) string {
	pad := length - len(text)
	if pad > 0 {
		return text + strings.Repeat(" ", pad)
	}

	if text[length-2] == ' ' {
		return text[:length-2] + "… "
	}

	return text[:length-1] + "…"
}

// still doesn't work for emojies
// line := strings.Split(text, "\n")[0]
// lineLen := runewidth.StringWidth(line)

// if lineLen <= length {
// 	return line + strings.Repeat(" ", length-lineLen)
// }

// runewidth.Truncate(line, length, "")

// if line[length-2] == ' ' {
// 	return runewidth.Truncate(line, length-2, "… ")
// }

// return runewidth.Truncate(line, length-1, "… ")
