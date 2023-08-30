package entity

import "time"

type Tokens struct {
	AccessToken 
	RefreshToken
}

type RefreshToken struct {
	Token 	  string    `json:"refresh_token"`
	ExpiresAt time.Time `json:"expires_at_refresh"`
}

type AccessToken struct {
	Token 	  string    `json:"access_token"`
	ExpiresAt time.Time `json:"expires_at_access"`
}