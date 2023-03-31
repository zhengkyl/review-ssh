package account

import (
	"bytes"
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type loginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func postAuth(client *retryablehttp.Client, loginData loginData) tea.Cmd {
	return func() tea.Msg {

		bsLoginData, err := json.Marshal(loginData)

		if err != nil {
			return common.AuthState{
				Authed: false,
			}
		}

		resp, err := client.Post("https://review-api.fly.dev/auth", "application/json", bytes.NewBuffer(bsLoginData))

		if err != nil {
			return common.AuthState{
				Authed: false,
			}
		}

		if resp.StatusCode != 200 {
			return common.AuthState{
				Authed: false,
			}
		}

		cookie := resp.Header.Get("Set-Cookie")

		var user common.User

		err = json.NewDecoder(resp.Body).Decode(&user)

		if err != nil {
			return common.AuthState{
				Authed: false,
			}
		}

		return common.AuthState{
			Authed: true,
			Cookie: cookie,
			User:   user,
		}
	}
}
