package server

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}