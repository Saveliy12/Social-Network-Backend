package database

import (
	"fmt"

	"github.com/Saveliy12/prod2/internal/models"
	"github.com/jmoiron/sqlx"
)

// UserRepositoryInterface определяет методы для работы с пользователями в базе данных
type UserRepositoryInterface interface {
	CreateUser(user models.RegistrationUser) (models.User, error)
	IsUnique(login, email, phone string) error
	GetUserByLogin(login string) (models.User, error)
	SetSession(userID uint, session models.Session) error
}

// UserRepository предоставляет реализацию UserRepositoryInterface
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository создает новый экземпляр UserRepository
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser добавляет нового пользователя в базу данных
func (s *UserRepository) CreateUser(user models.RegistrationUser) (models.User, error) {
	query := `
        INSERT INTO users (login, email, phone, password, createdat)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, login, email, phone
    `

	var newUser models.User
	err := s.db.QueryRow(query, user.Login, user.Email, user.Phone, user.Password, user.CreatedAt).Scan(
		&newUser.ID, &newUser.Login, &newUser.Email, &newUser.Phone,
	)
	if err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

// IsUnique проверяет, что указанные значения login, email и phone уникальны в базе данных
func (s *UserRepository) IsUnique(login, email, phone string) error {
	// Combining the queries into a single query
	query := `
		SELECT 
			(SELECT COUNT(*) FROM users WHERE login = $1) AS loginCount,
			(SELECT COUNT(*) FROM users WHERE email = $2) AS emailCount,
			(SELECT COUNT(*) FROM users WHERE phone = $3) AS phoneCount
	`
	var counts struct {
		LoginCount int `db:"loginCount"`
		EmailCount int `db:"emailCount"`
		PhoneCount int `db:"phoneCount"`
	}

	if err := s.db.Get(&counts, query, login, email, phone); err != nil {
		return fmt.Errorf("failed to check uniqueness in the database: %v", err)
	}

	if counts.LoginCount > 0 {
		return fmt.Errorf("login already exists")
	}
	if counts.EmailCount > 0 {
		return fmt.Errorf("email already exists")
	}
	if counts.PhoneCount > 0 {
		return fmt.Errorf("phone number already exists")
	}

	return nil
}

func (s *UserRepository) GetUserByLogin(login string) (models.User, error) {
	var user models.User
	query := "SELECT id, login, email, phone, password FROM users WHERE login = $1"
	err := s.db.Get(&user, query, login)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user by login: %v", err)
	}
	return user, nil
}

func (s *UserRepository) SetSession(userID uint, session models.Session) error {
	_, err := s.db.Exec("UPDATE sessions SET refresh_token = ?, expires_at = ? WHERE userID = ?", session.RefreshToken, session.ExpiresAt, userID)
	return err
}
