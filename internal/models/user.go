package models

import "time"

type RegistrationUser struct {
	Login     string    `json:"login" db:"login"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"createdAt" db:"createdAt"`
}

type LoginUser struct {
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type User struct {
	ID       uint   `json:"id" db:"id"`
	Login    string `json:"login" db:"login"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Password string `json:"-"`
}

type Session struct {
	RefreshToken string    `json:"refreshToken" db:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expiresAt"`
}
