package usecase

import (
	"context"
	"goauth/internal/entity"
	"goauth/pkg/manager"
	"goauth/pkg/hash"
	"log"
	"time"
)

type UserStorage interface {
	CreateUser(context.Context, *entity.UserInput, *entity.RefreshToken) error
	GetUser(context.Context, string, string) (string, error)
	UpdateUser(context.Context, string, *entity.RefreshToken) error
	Refresh(context.Context, string, *entity.RefreshToken) error
	GetUsername(context.Context, string) (string, time.Time, error)
}

type AuthInteractor struct {
	userStorage    UserStorage
	tokenManager   manager.TokenManager
	passwordHasher hash.PasswordHasher
}

func NewAuthInteractor(userStorage UserStorage, tokenManager manager.TokenManager, hasher hash.PasswordHasher) *AuthInteractor {
	return &AuthInteractor{
		userStorage:    userStorage,
		tokenManager:   tokenManager,
		passwordHasher: hasher,
	}
}

func (t *AuthInteractor) SignUp(ctx context.Context, user *entity.UserInput) (string, string, error) {
	user.Password = t.passwordHasher.Hash(user.Password)
	
	refreshToken, err := t.tokenManager.NewRefreshToken()
	if err != nil {
		log.Println(err)
        return "", "", err
    }

	accessToken, err := t.tokenManager.NewJWT(user.Username)
	if err != nil {
        return "", "", err
    }

	return accessToken, refreshToken.Token, t.userStorage.CreateUser(ctx, user, refreshToken)
}

func (t *AuthInteractor) SignIn(ctx context.Context, user *entity.UserInput) (string, string, error) {
	user.Password = t.passwordHasher.Hash(user.Password)

	username, err := t.userStorage.GetUser(ctx, user.Username, user.Password)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := t.tokenManager.NewRefreshToken()
	if err != nil {
        return "", "", err
    }

	err = t.userStorage.UpdateUser(ctx, username, refreshToken)

	accessToken, err := t.tokenManager.NewJWT(user.Username)
	if err != nil {
        return "", "", err
    }

	return accessToken, refreshToken.Token, nil
}

func (t *AuthInteractor) RefreshToken(ctx context.Context, oldToken string) (string, string, error) {
	newRefreshToken, err := t.tokenManager.NewRefreshToken()
	if err!= nil {
        return "", "", err
    }

	username, _, err := t.userStorage.GetUsername(ctx, oldToken)
	if err != nil {
		return "", "", err
	}

	err = t.userStorage.Refresh(ctx, oldToken, newRefreshToken)
	if err != nil {
		return "", "", err
	}

	accessToken, err := t.tokenManager.NewJWT(username)
	if err != nil {
        return "", "", err
    }

	return accessToken, newRefreshToken.Token, nil
}
