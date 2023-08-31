package server

import "goauth/internal/entity"

type errorMessage struct {
	Message string `json:"message"`
}

type outputTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func successResponse(t *entity.Tokens) outputTokens {
	return outputTokens{
        AccessToken:  t.AccessToken,
        RefreshToken: t.RefreshToken,
    }
}

func errorResponse(message string) errorMessage {
	return errorMessage{
        Message: message,
    }
}