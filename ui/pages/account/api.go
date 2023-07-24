package account

import (
	"bytes"
	"encoding/json"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/common"
)

type signUpData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signUpRes struct {
	ok  bool
	err string
}

func postSignUp(client *retryablehttp.Client, data signUpData) tea.Msg {
	bsLoginData, err := json.Marshal(data)

	if err != nil {
		return signUpRes{false, err.Error()}
	}

	resp, err := client.Post(common.ReviewBase+"/users", "application/json", bytes.NewBuffer(bsLoginData))

	if err != nil {
		return signUpRes{false, err.Error()}
	}

	if resp.StatusCode != 200 {
		return signUpRes{false, "Email already registered."}
	}

	var user common.User

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return signUpRes{false, err.Error()}
	}

	return signUpRes{true, err.Error()}
}

type signInData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type signInRes struct {
	ok  bool
	err string
}

func postSignIn(client *retryablehttp.Client, data signInData) tea.Msg {
	bsLoginData, err := json.Marshal(data)

	if err != nil {
		return signInRes{false, err.Error()}
	}

	resp, err := client.Post(common.ReviewBase+"/auth", "application/json", bytes.NewBuffer(bsLoginData))

	if err != nil {
		return signInRes{false, err.Error()}
	}

	if resp.StatusCode != 200 {
		return signInRes{false, "Wrong email or password."}
	}

	var cookie string
	for _, c := range resp.Cookies() {
		if c.Name == "id" {
			cookie = c.Value
		}
	}

	var user common.User

	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		return signInRes{false, err.Error()}
	}

	return common.AuthState{
		Authed: true,
		Cookie: cookie,
		User:   user,
	}
}
