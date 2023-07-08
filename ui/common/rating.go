package common

// ▣ ☑ ☐ 🗹 ▢ ⬚ ⛶ ▢
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
)

func RenderRating(before, during, after bool) string {
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
