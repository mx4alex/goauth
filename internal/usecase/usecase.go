package usecase

import (
	"context"
	"goauth/internal/entity"
	"goauth/pkg/manager"
	"goauth/pkg/hash"
	"errors"
	"log"
	"time"
)

type UserStorage interface {
	CreateUser(context.Context, *entity.UserSignUp, *entity.RefreshToken) error
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

func (t *AuthInteractor) SignUp(ctx context.Context, user *entity.UserSignUp) (*entity.Tokens, error) {
	user.Password = t.passwordHasher.Hash(user.Password)

	_, err := t.userStorage.GetUser(ctx, user.Username, user.Password)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	
	refreshToken, err := t.tokenManager.NewRefreshToken()
	if err != nil {
		log.Println(err)
        return nil, err
    }

	accessToken, err := t.tokenManager.NewJWT(user.Username)
	if err != nil {
        return nil, err
    }

	tokens := &entity.Tokens{
		AccessToken: accessToken,
		RefreshToken: refreshToken.Token,
	}

	err = t.userStorage.CreateUser(ctx, user, refreshToken)
	if err != nil {
        return nil, err
    }

	return tokens, nil
}

func (t *AuthInteractor) SignIn(ctx context.Context, user *entity.UserSignIn) (*entity.Tokens, error) {
	user.Password = t.passwordHasher.Hash(user.Password)

	username, err := t.userStorage.GetUser(ctx, user.Username, user.Password)
	if err != nil {
		return nil, err
	}

	refreshToken, err := t.tokenManager.NewRefreshToken()
	if err != nil {
        return nil, err
    }

	err = t.userStorage.UpdateUser(ctx, username, refreshToken)

	accessToken, err := t.tokenManager.NewJWT(user.Username)
	if err != nil {
        return nil, err
    }

	tokens := &entity.Tokens{
		AccessToken: accessToken,
		RefreshToken: refreshToken.Token,
	}

	return tokens, nil
}

func (t *AuthInteractor) RefreshToken(ctx context.Context, oldToken string) (*entity.Tokens, error) {
	newRefreshToken, err := t.tokenManager.NewRefreshToken()
	if err != nil {
        return nil, err
    }

	username, _, err := t.userStorage.GetUsername(ctx, oldToken)
	if err != nil {
		return nil, err
	}

	err = t.userStorage.Refresh(ctx, oldToken, newRefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := t.tokenManager.NewJWT(username)
	if err != nil {
        return nil, err
    }

	tokens := &entity.Tokens{
		AccessToken: accessToken,
		RefreshToken: newRefreshToken.Token,
	}

	return tokens, nil
}
