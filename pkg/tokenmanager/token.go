package tokenmanager

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Saveliy12/prod2/pkg/logger"
	"github.com/golang-jwt/jwt"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type TokenManagerInterface interface {
	NewJWT(userId string, ttl time.Duration) (string, error)
	ParseJWT(accessToken string) (uint, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKey string
	log        logger.LoggerInterface
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey, log: logger.GetLogger()}, nil
}

func (m *Manager) NewJWT(userID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(ttl).Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) ParseJWT(accessToken string) (uint, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(uint), nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
