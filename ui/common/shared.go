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
	Id         int32
	Name       string
	Email      string
	Created_at time.Time
	Updated_at time.Time
}
