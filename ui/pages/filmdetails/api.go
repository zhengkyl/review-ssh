package filmdetails

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zhengkyl/review-ssh/ui/common"
)

const patchEndpoint = common.ReviewBase + "/reviews/Film/"

func patchReviewCmd(g common.Global, tmdb_id int, updates map[string]interface{}) tea.Cmd {
	return common.Fetch[common.Review](g, "PATCH", patchEndpoint+strconv.Itoa(tmdb_id), updates, func(data common.Review, err error) tea.Msg {
		if err == nil {
			g.ReviewMap[tmdb_id] = data
		}
		return nil
	})
}

const postEndpoint = common.ReviewBase + "/reviews"

func postReviewCmd(g common.Global, tmdb_id int, status string) tea.Cmd {
	data := map[string]interface{}{
		"tmdb_id":  tmdb_id,
		"category": "Film",
		"status":   status,
	}
	return common.Fetch[common.Review](g, "POST", postEndpoint, data, func(data common.Review, err error) tea.Msg {
		if err == nil {
			g.ReviewMap[tmdb_id] = data
		}
		return nil
	})
}
