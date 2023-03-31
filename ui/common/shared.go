package common

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type Shared struct {
	AuthState  AuthState
	HttpClient retryablehttp.Client
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
