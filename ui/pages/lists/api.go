package lists

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

// TODO use pagination, but for now 50 is more than enough
const reviewsEndpoint = "https://review-api.fly.dev/reviews?category=Film&per_page=50"

func getReviewsCmd(client *retryablehttp.Client, user_id int) tea.Cmd {
	if user_id == common.GuestAuthState.User.Id {
		return common.GetCmd[common.PageResult[common.Review]](client, reviewsEndpoint)
	}

	return common.GetCmd[common.PageResult[common.Review]](client, reviewsEndpoint+"?user_id="+strconv.Itoa(user_id))
}
