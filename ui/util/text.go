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
