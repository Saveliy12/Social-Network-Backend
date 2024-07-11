package api

import (
	"net/http"
	"time"

	"github.com/Saveliy12/prod2/internal/models"
	"github.com/Saveliy12/prod2/internal/service"
	"github.com/Saveliy12/prod2/pkg/logger"
	"github.com/gin-gonic/gin"
)

// AuthHandlerInterface определяет методы для работы с аутентификацией
type AuthHandlerInterface interface {
	RegisterUserHandler(c *gin.Context)
	SignInUserHandler(c *gin.Context)
}

// AuthHandler предоставляет обработчики для аутентификации
type AuthHandler struct {
	authService service.AuthService
	log         logger.LoggerInterface
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         logger.GetLogger(),
	}
}

// RegisterUserHandler обрабатывает запрос на регистрацию нового пользователя
func (a *AuthHandler) RegisterUserHandler(c *gin.Context) {
	var newUser models.RegistrationUser
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, что все необходимые данные получены
	if newUser.Login == "" || newUser.Email == "" || newUser.Phone == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Устанавливаем текущее время в поле CreatedAt
	newUser.CreatedAt = time.Now()

	// Передаем нового пользователя в сервис для создания
	createdUser, err := a.authService.RegisterUser(newUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register user: " + err.Error()})
		return
	}

	// Возвращаем успешный ответ с созданным пользователем
	c.JSON(http.StatusCreated, createdUser)
}

func (a *AuthHandler) LoginUserHandler(c *gin.Context) {
	var credentials models.LoginUser
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := a.authService.AuthenticateUser(credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login or password"})
		return
	}

	tokens, err := a.authService.AuthorizeUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to  user: " + err.Error()})
		return
	}

	tokenResponse := struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	c.JSON(http.StatusOK, tokenResponse)
}

func (a *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccessToken, err := a.authService.NewRefreshToken(requestBody.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccessToken})
}
