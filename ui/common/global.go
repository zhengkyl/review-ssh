package common

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/zhengkyl/review-ssh/ui/keymap"
	"github.com/zhengkyl/review-ssh/ui/styles"
)

type Global struct {
	AuthState  AuthState
	HttpClient retryablehttp.Client
	Styles     *styles.Styles
	KeyMap     *keymap.KeyMap
}

type AuthState struct {
	Authed bool
	Cookie string
	User   User
}

type User struct {
	Id         int32     `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

var GuestAuthState = AuthState{
	Authed: true,
	Cookie: "guestcookie",
	User: User{
		Id:    -1,
		Name:  "Guest",
		Email: "guest@zhengkyl.com",
	},
}
