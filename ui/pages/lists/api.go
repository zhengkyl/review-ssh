package lists

import (
	"encoding/json"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type res struct {
	ok  bool
	err string
}

type getResponse struct {
	Results []common.Review
	// unused fields
	// Page          int
	// Total_Pages   int
	// Total_Results int
}

type params map[string]string

// TODO use pagination, but for now 50 is more than enough
const reviewsEndpoint = "https://review-api.fly.dev/reviews?per_page=50"

func getReviewsCmd(client *retryablehttp.Client, user_id int) tea.Cmd {
	return func() tea.Msg {
		if user_id == common.GuestAuthState.User.Id {
			return getReviews(client, params{})
		}
		return getReviews(client, params{"user_id": strconv.Itoa(user_id)})
	}
}

func getReviews(client *retryablehttp.Client, params params) tea.Msg {
	endpoint := reviewsEndpoint

	if len(params) > 0 {
		endpoint += "?"

		for key, value := range params {
			endpoint += key + "=" + value + "&"
		}
		endpoint = endpoint[:len(endpoint)-1]
	}

	resp, err := client.Get(endpoint)

	if err != nil {
		return res{false, err.Error()}
	}

	if resp.StatusCode != 200 {
		return res{false, "Something went wrong."}
	}

	var response getResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return res{false, err.Error()}
	}

	return response.Results
}
