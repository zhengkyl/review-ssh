package common

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/keymap"
)

type Global struct {
	AuthState  *AuthState
	Config     Config
	HttpClient *retryablehttp.Client
	KeyMap     *keymap.KeyMap
}

type Config struct {
	TMDB_API_KEY string
}

type AuthState struct {
	Authed bool
	Cookie string
	User   User
}

var GuestAuthState = AuthState{
	Authed: true,
	Cookie: "guestcookie",
	User: User{
		Id:   -1,
		Name: "Guest",
	},
}
