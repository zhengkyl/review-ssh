package filmdetails

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
	"github.com/zhengkyl/review-ssh/ui/common/enums"
)

const patchEndpoint = "https://review-api.fly.dev/reviews/Film/"

func patchReviewCmd(client *retryablehttp.Client, tmdb_id int, updates map[string]string) tea.Cmd {
	return common.Fetch[common.Review](client, "PATCH", patchEndpoint+strconv.Itoa(tmdb_id), updates, func(data common.Review, err error) tea.Msg {
		return nil
	})
}

const postEndpoint = "https://review-api.fly.dev/reviews"

func postReviewCmd(client *retryablehttp.Client, tmdb_id int, status enums.Status) tea.Cmd {
	data := map[string]string{
		"tmdb_id":  strconv.Itoa(tmdb_id),
		"category": "Film",
		"status":   status.String(),
	}
	return common.Fetch[common.Review](client, "POST", postEndpoint+strconv.Itoa(tmdb_id), data, func(data common.Review, err error) tea.Msg { return nil })
}
