package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Saveliy12/prod2/internal/models"
)

func ValidateUser(user models.RegistrationUser) error {

	if err := validateLogin(user.Login); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := validateEmail(user.Email); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := validatePassword(user.Password); err != nil {
		return fmt.Errorf("%w", err)
	}

	if err := validatePhoneNumber(user.Phone); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func validateLogin(login string) error {
	if len(login) > 30 {
		return fmt.Errorf("max login length is 30 characters")
	}

	// Проверка на соответствие шаблону [a-zA-Z0-9-]+
	pattern := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	if !pattern.MatchString(login) {
		return fmt.Errorf("the login can contain only Latin letters and numbers")
	}
	return nil
}

func validateEmail(email string) error {
	// Проверка минимальной длины
	if len(email) < 1 {
		return fmt.Errorf("min email length is 1 characters")
	}

	// Проверка максимальной длины
	if len(email) > 50 {
		return fmt.Errorf("max email length is 50 characters")
	}

	// Проверка на наличие символа @
	if !strings.Contains(email, "@") {
		return fmt.Errorf("the email must contain @")
	}

	return nil
}

func validatePassword(password string) error {
	// Проверка минимальной длины
	if len(password) < 6 {
		return fmt.Errorf("min password length is 6 characters")
	}

	// Проверка максимальной длины
	if len(password) > 100 {
		return fmt.Errorf("max password length is 100 characters")
	}

	// Проверка на наличие латинских символов в верхнем и нижнем регистрах
	if matched, _ := regexp.MatchString(`[a-z]+`, password); !matched {
		return fmt.Errorf("password must contain at least one lowercase Latin character")
	}
	if matched, _ := regexp.MatchString(`[A-Z]+`, password); !matched {
		return fmt.Errorf("password must contain at least one uppercase Latin character")
	}

	// Проверка на наличие хотя бы одной цифры
	if matched, _ := regexp.MatchString(`[0-9]+`, password); !matched {
		return fmt.Errorf("password must contain at least one digit")
	}

	// Проверка на наличие других специальных символов
	if matched, _ := regexp.MatchString(`[!@#$%^&*()-_+=~{}[\]|;:/?,.<>]+`, password); !matched {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func validatePhoneNumber(phone string) error {
	// Проверка максимальной длины
	if len(phone) > 20 {
		return fmt.Errorf("max phone number lenght is 20 characters")
	}

	// Проверка на соответствие шаблону \+[\d]+
	pattern := regexp.MustCompile(`^\+[\d]+$`)
	if !pattern.MatchString(phone) {
		return fmt.Errorf("invalid phone number format")
	}

	return nil
}
