package usecase

import (
	"goauth/internal/entity"
	"crypto/sha1"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
	"fmt"
	"context"
)

type UserStorage interface {
	CreateUser(context.Context, *entity.User) error
	GetUser(context.Context, string, string) (*entity.User, error)
}

type AuthInteractor struct {
	userStorage UserStorage

	hashSalt 	   string
	expireDuration time.Duration
	signingKey     []byte
}

func NewAuthInteractor(userStorage UserStorage) *AuthInteractor {
	return &AuthInteractor{
		userStorage: 	userStorage,
		hashSalt:    	"<PASSWORD>",
        expireDuration: time.Hour * 24,
        signingKey: 	[]byte("secret"),
	}
}

func (t *AuthInteractor) SignUp(ctx context.Context, user *entity.User) error {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(t.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))
	return t.userStorage.CreateUser(ctx, user)
}

func (t *AuthInteractor) SignIn(ctx context.Context, user *entity.User) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(t.hashSalt))
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := t.userStorage.GetUser(ctx, user.Username, user.Password)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(t.expireDuration)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: user.Username,
	})

	return token.SignedString(t.signingKey)
}
