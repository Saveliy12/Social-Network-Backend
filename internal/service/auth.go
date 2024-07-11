package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/Saveliy12/prod2/internal/database"
	"github.com/Saveliy12/prod2/internal/models"
	"github.com/Saveliy12/prod2/internal/utils"
	"github.com/Saveliy12/prod2/pkg/hash"
	tokenmanager "github.com/Saveliy12/prod2/pkg/tokenmanager"
	"golang.org/x/crypto/bcrypt"
)

// AuthServiceInterface определяет методы для работы с аутентификацией
type AuthServiceInterface interface {
	RegisterUser(newUser models.RegistrationUser) (models.User, error)
	AuthenticateUser(credentials models.LoginUser) (models.User, error)
	AuthorizeUserUser(userID uint) (string, error)
	NewRefreshToken(refreshToken string) (string, error)
}

// AuthService предоставляет реализацию AuthServiceInterface
type AuthService struct {
	tokenManager   tokenmanager.TokenManagerInterface
	userRepository database.UserRepositoryInterface

	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(tokenManager tokenmanager.TokenManagerInterface, userRepository database.UserRepositoryInterface,
	accessTokenTTL time.Duration, refreshTokenTTL time.Duration) *AuthService {
	return &AuthService{
		tokenManager:    tokenManager,
		userRepository:  userRepository,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *AuthService) RegisterUser(newUser models.RegistrationUser) (models.User, error) {
	if err := s.userRepository.IsUnique(newUser.Login, newUser.Email, newUser.Phone); err != nil {
		return models.User{}, err
	}

	if err := utils.ValidateUser(newUser); err != nil {
		return models.User{}, err
	}

	hashedPassword, err := hash.HashPassword(newUser.Password)
	newUser.Password = hashedPassword

	return s.userRepository.CreateUser(newUser)
}

func (s *AuthService) AuthenticateUser(credentials models.LoginUser) (uint, error) {
	user, err := s.userRepository.GetUserByLogin(credentials.Login)
	if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return 0, errors.New("invalid login or password")
	}

	return user.ID, nil
}

func (s *AuthService) AuthorizeUser(userID uint) (tokenmanager.Tokens, error) {
	var (
		res tokenmanager.Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(strconv.FormatUint(uint64(userID), 10), s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := models.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.userRepository.SetSession(userID, session)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (s *AuthService) NewRefreshToken(refreshToken string) (string, error) {
	userID, err := s.tokenManager.ParseJWT(refreshToken)
	if err != nil {
		return "", err
	}

	userID_string := strconv.FormatUint(uint64(userID), 10)

	newAccessToken, err := s.tokenManager.NewJWT(userID_string, s.accessTokenTTL)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}
