package manager

import (
	"fmt"
	"math/rand"
	"time"
	"goauth/internal/entity"
	"github.com/dgrijalva/jwt-go/v4"
)

type TokenManager interface {
	NewJWT(username string) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (*entity.RefreshToken, error)
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(signingKey string, accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *Manager {
	return &Manager{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (m *Manager) NewJWT(username string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &entity.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(m.accessTokenTTL)),
			IssuedAt:  jwt.At(time.Now()),
		},
		Username: username,
	})

	return accessToken.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken() (*entity.RefreshToken, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return nil, err
	}

	token := new(entity.RefreshToken)
	token.Token = fmt.Sprintf("%x", b)
	token.ExpiresAt = time.Now().Add(m.refreshTokenTTL)

	return token, nil
}