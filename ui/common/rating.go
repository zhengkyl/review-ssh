package common

var (
	ratings = []string{
		"░░ ░░",
		"██ ░░",
		"░░ ██",
		"██ ██",
	}
)

func RenderRating(before, during, after bool) string {
	ratingIndex := 0
	if during {
		ratingIndex += 1
	}
	if after {
		ratingIndex += 2
	}
	return ratings[ratingIndex]
}
